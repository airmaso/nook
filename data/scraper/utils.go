package scraper

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"

	"nook/core/models"
)

const (
	// All-time rankings leaderboard page
	GlobalLeaderboardURL = "https://starbattle.puzzlebaron.com/halloffame.php"

	// Monthly competition leaderboard page
	MonthlyLeaderboardURL = "https://starbattle.puzzlebaron.com/contest1.php"

	// Lifetime player solving statistics page
	PlayerProfileURL = "https://starbattle.puzzlebaron.com/profile.php?u=%s"
)

// Extracts and returns the text content from a leaderboard table.
// Preserves rank order by processing the leaderboard iteratively (top-down).
//
// If the leaderboard is monthly, each row is structured as:
//
//	[]string{
//		uid,
//		monthlyRank,
//		username,
//		lastActive,
//		monthlySolvedFraction,
//		monthlyRate,
//		monthlyAvgTime,
//		monthlyAvgPoints,
//		monthlyTotalPoints,
//	}
//
// If the leaderboard is global, each row is structured as:
//
//	[]string{
//		uid,
//		globalRank,
//		username,
//		lastGame,
//		lifetimePoints,
//	}
func GetLeaderboardTextDataV1(lb models.Leaderboard) ([][]string, error) {
	url := MonthlyLeaderboardURL
	if lb == models.GlobalLeaderboard {
		url = GlobalLeaderboardURL
	}

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Fetching leaderboard: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status: %d", response.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failure parsing response body")
	}

	dataRows := doc.Find("table.winners tbody tr")
	if dataRows == nil {
		return nil, fmt.Errorf("No data rows found")
	}

	var textData [][]string

	// Iterate over each row of the leaderboard table
	dataRows.Each(func(rowIndex int, row *goquery.Selection) {
		// Skip the header row
		if rowIndex == 0 {
			return
		}

		rowData := []string{}

		// Iterate over each column of this row
		row.Find("td").Each(func(colIndex int, td *goquery.Selection) {
			// Append each text column
			rowData = append(rowData, td.Text())

			// Extract the UID
			if link := td.Find("a"); link.Length() != 0 {
				if href, exists := link.Attr("href"); exists {
					if _, uid, found := strings.Cut(href, "?u="); found {
						rowData = append([]string{uid}, rowData...)
					}
				}
				// fmt.Printf("[@%v] %v\n", link.Text(), href)
			}
		})

		// Append the row data to the bulk table data
		textData = append(textData, rowData)
	})

	return textData, nil
}

// Extracts and returns the text content from a leaderboard table.
// Improves upon [GetLeaderboardTextDataV1] by processing rows concurrently
// with goroutines, writing each row to a preallocated slice to preserve
// rank order without explicit synchronization overhead.
//
// If the leaderboard is monthly, each row is structured as:
//
//	[]string{
//	  uid,
//	  monthlyRank,
//	  lastActive,
//	  monthlySolvedFraction,
//	  monthlyRate,
//	  monthlyAvgTime,
//	  monthlyAvgPoints,
//	  monthlyTotalPoints,
//	}
//
// If the leaderboard is global, each row is structured as:
//
//	[]string{
//	  uid,
//	  globalRank,
//	  lastGame,
//	  lifetimePoints,
//	}
func GetLeaderboardTextDataV2(lb models.Leaderboard) ([][]string, error) {
	url := MonthlyLeaderboardURL
	if lb == models.GlobalLeaderboard {
		url = GlobalLeaderboardURL
	}

	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Fetching leaderboard: %w", err)
	}
	defer response.Body.Close()

	// Parse response body
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failure parsing response body")
	}

	// Find all leaderboard table rows
	rows := doc.Find("table.winners tbody tr")

	// Preallocated so each goroutine can write directly to textData[rowIndex-1],
	// preservering rank order without relying on completion order
	// No mutex/channel needed: concurrent writes to distinct slice indices are
	// race-free since each goroutine only touches its own memory location (&textData[rowIndex-1])
	textData := make([][]string, rows.Length()-1)
	var wg sync.WaitGroup

	// Concurrently process each row
	rows.Each(func(rowIndex int, row *goquery.Selection) {
		// Skip header row
		if rowIndex == 0 {
			return
		}

		// Process this row in it's own goroutine
		wg.Go(func() {
			// Extract and collect the column text from this row
			var cols []string
			row.Find("td").Each(func(colIndex int, td *goquery.Selection) {
				cols = append(cols, strings.TrimSpace(td.Text()))

				// Extract the uid
				if link := td.Find("a"); link.Length() != 0 {
					if href, exists := link.Attr("href"); exists {
						if _, uid, found := strings.Cut(href, "?u="); found {
							cols = append([]string{uid}, cols...)
						}
					}
					// fmt.Printf("[@%v] %v\n", link.Text(), href)
				}
			})

			// Each goroutine writes to its own index
			textData[rowIndex-1] = cols
		})
	})

	// Block until all goroutines finish
	wg.Wait()

	return textData, nil
}
