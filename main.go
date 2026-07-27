package main

import (
	"fmt"

	"nook/core/models"
	"nook/data/scraper"
)

func main() {
	for _, lb := range []models.Leaderboard{
		models.MonthlyLeaderboard,
		models.GlobalLeaderboard,
	} {
		users, stats, err := scraper.GetLeaderboardUsers(lb)
		if err != nil {
			panic("couldn't get leaderboard users")
		}

		fmt.Printf("%s Took: %v\n\n", lb.String(), stats.ScrapeTime)

		for _, u := range users {
			fmt.Println(u.DetailedString())
		}

		fmt.Println()
	}
}
