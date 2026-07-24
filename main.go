package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"nook/core/models"
	"nook/data/scraper"
)

func main() {
	jobs := []struct {
		leaderboard models.Leaderboard
		scrape      func(models.Leaderboard) ([][]string, error)
	}{
		{models.MonthlyLeaderboard, scraper.GetLeaderboardTextDataV1},
		{models.MonthlyLeaderboard, scraper.GetLeaderboardTextDataV2},

		{models.GlobalLeaderboard, scraper.GetLeaderboardTextDataV1},
		{models.GlobalLeaderboard, scraper.GetLeaderboardTextDataV2},
	}

	ratio := math.Pow(10, 2) // 2-decimal place precision

	for i := 0; i < len(jobs); i += 2 {
		leaderboard := jobs[i].leaderboard
		oldScrape := jobs[i].scrape
		newScrape := jobs[i+1].scrape

		// Old scrape benchmark
		start := time.Now()
		oldScrape(leaderboard)
		oldScrapeTime := time.Since(start)

		// New scrape benchmark
		start = time.Now()
		data, err := newScrape(leaderboard)
		if err != nil {
			panic("newScrape failed")
		}

		newScrapeTime := time.Since(start)

		diff := oldScrapeTime - newScrapeTime
		mult := float64(oldScrapeTime.Milliseconds()) / float64(newScrapeTime.Milliseconds())

		fmt.Printf("---- %v ----\n", leaderboard.String())
		fmt.Printf("Old scrape: %vms\n", oldScrapeTime.Milliseconds())
		fmt.Printf("New scrape: %vms\n", newScrapeTime.Milliseconds())
		fmt.Printf("  + Improvement: %vms\n", diff.Milliseconds())
		fmt.Printf("  + Throughput: %vx\n", math.Round(mult*ratio)/ratio)

		fmt.Println()

		for _, row := range data {
			fmt.Println(strings.Join(row, "|"))
		}

		fmt.Println()
	}
}
