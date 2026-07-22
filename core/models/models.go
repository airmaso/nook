package models

import (
	"fmt"
	"time"
)

type Leaderboard int

const (
	MonthlyLeaderboard Leaderboard = iota
	GlobalLeaderboard
)

type User struct {
	UID            int
	Rank           int
	Username       string
	LastActive     time.Time
	Solved         int
	Attempted      int
	AcceptanceRate float64
	AverageTime    float64
	AveragePoints  float64
	TotalPoints    int
}

func (u *User) String() string {
	rank := fmt.Sprintf("%3d", u.Rank)                                        // right-align in 3 chars; '%d' right-aligns by default
	username := fmt.Sprintf("%-33s", "@"+u.Username)                          // left-align in 33 chars; '-' = left
	totalPoints := fmt.Sprintf("%-15s", fmt.Sprintf("%d pts", u.TotalPoints)) // left-align in 15 chars

	return fmt.Sprintf("%s %s %s", rank, username, totalPoints)
}
