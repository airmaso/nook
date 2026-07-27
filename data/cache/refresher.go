package cache

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"nook/core/models"
)

// Reads REFRESH_INTERVAL (seconds) from the environment.
//
// Falls back to a default of 2 minutes if the parse fails.
func ParseRefreshInterval() time.Duration {
	const defaultInterval = 2 * time.Minute

	raw := os.Getenv("REFRESH_INTERVAL")
	if raw == "" {
		return defaultInterval
	}

	seconds, err := strconv.Atoi(raw)
	if err != nil || seconds <= 0 {
		slog.Warn("invalid REFRESH_INTERVAL, using default",
			"value", raw, "default", defaultInterval,
		)
		return defaultInterval
	}

	return time.Duration(seconds) * time.Second
}

// Launches a background goroutine loop that refreshes each leaderboard
// on a fixed interval until the context is cancelled
func (cache *LeaderboardCache) LaunchRefresher(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Prime the cache immediately on launch
		cache.refreshAll()

		for {
			select {
			case <-ctx.Done():
				slog.Info("refresher halted")
				return
			case <-ticker.C:
				cache.refreshAll()
			}
		}
	}()
}

// Refreshes both leaderboards
func (cache *LeaderboardCache) refreshAll() error {
	var wg sync.WaitGroup

	leaderboards := []models.Leaderboard{models.MonthlyLeaderboard, models.GlobalLeaderboard}

	// Refresh both leaderboards concurrently
	for _, lb := range leaderboards {
		wg.Go(func() {
			if err := cache.Refresh(lb); err != nil {
				slog.Error("refresh failed", "leaderboard", "lb", "error", err)
			}

			// slog.Info("successfuly refreshed", "leaderboard", lb.String())
		})
	}

	wg.Wait()
	return nil
}
