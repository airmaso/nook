package main

import (
	"fmt"
	"time"

	"nook/core/models"
)

func main() {
	airmaso := models.User{
		UID:            1,
		Rank:           1,
		Username:       "airmaso",
		LastActive:     time.Now(),
		Solved:         999,
		Attempted:      999,
		AcceptanceRate: float64(100),
		AverageTime:    float64(99),
		AveragePoints:  float64(555),
		TotalPoints:    1_000_000,
	}

	fmt.Println(airmaso.String())
}
