package cache

import (
	"fmt"
	"sync"

	"nook/core/models"
	"nook/data/scraper"
)

// A cached leaderboard
type CachedLeaderboard struct {
	prev  []*models.User
	curr  []*models.User
	stats *scraper.ScrapeStats
}

// The main interface for the leaderboard cache
type LeaderboardCache struct {
	mu                 sync.RWMutex
	monthlyLeaderboard CachedLeaderboard
	globalLeaderboard  CachedLeaderboard
}

// Helper that returns a reference to a cached leaderboard
func getCachedLeaderboard(cache *LeaderboardCache, lb models.Leaderboard) (*CachedLeaderboard, error) {
	switch lb {
	case models.MonthlyLeaderboard:
		return &cache.monthlyLeaderboard, nil
	case models.GlobalLeaderboard:
		return &cache.globalLeaderboard, nil
	default:
		return &CachedLeaderboard{}, fmt.Errorf(
			"unknown leaderboard type %q", lb,
		)
	}
}

// Returns the previously fetched users of a leaderboard
func (cache *LeaderboardCache) GetPrevUsers(lb models.Leaderboard) ([]*models.User, error) {
	cachedLeaderboard, err := getCachedLeaderboard(cache, lb)
	if err != nil {
		return nil, err
	}

	return cachedLeaderboard.prev, nil
}

// Returns the most recently fetched users of a leaderboard
func (cache *LeaderboardCache) GetCurrUsers(lb models.Leaderboard) ([]*models.User, error) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	cachedLeaderboard, err := getCachedLeaderboard(cache, lb)
	if err != nil {
		return nil, err
	}

	return cachedLeaderboard.curr, nil
}

// Returns the previously fetched users of a leaderboard
func (cache *LeaderboardCache) GetPrevUser(lb models.Leaderboard, rank int) (*models.User, error) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	cachedLeaderboard, err := getCachedLeaderboard(cache, lb)
	if err != nil {
		return nil, err
	}

	if rank <= 0 || rank > len(cachedLeaderboard.curr) {
		return nil, fmt.Errorf(
			"rank %d not in range [1, %d]", rank, len(cachedLeaderboard.curr),
		)
	}

	return cachedLeaderboard.prev[rank-1], nil
}

// Returns a specific user from the most recently fetched users of a leaderboard
func (cache *LeaderboardCache) GetCurrUser(lb models.Leaderboard, rank int) (*models.User, error) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	cachedLeaderboard, err := getCachedLeaderboard(cache, lb)
	if err != nil {
		return nil, err
	}

	if rank <= 0 || rank > len(cachedLeaderboard.curr) {
		return nil, fmt.Errorf(
			"rank %d not in range [1, %d]", rank, len(cachedLeaderboard.curr),
		)
	}

	return cachedLeaderboard.curr[rank-1], nil
}

// Refreshes a cached leaderboard
func (cache *LeaderboardCache) Refresh(lb models.Leaderboard) error {
	// NOTE: Allow the scrape to occur outside the lock since
	// we don't want the network call holding a write lock (and blocking readers)
	users, stats, err := scraper.GetLeaderboardUsers(lb)
	if err != nil {
		return fmt.Errorf(
			"cache refresh failed",
		)
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()

	cachedLeaderboard, err := getCachedLeaderboard(cache, lb)
	if err != nil {
		return err
	}

	cachedLeaderboard.prev = cachedLeaderboard.curr
	cachedLeaderboard.curr = users
	cachedLeaderboard.stats = stats

	return nil
}
