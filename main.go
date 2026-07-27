package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"syscall"

	"nook/core/models"
	"nook/data/cache"
)

type nookServer struct {
	cache *cache.LeaderboardCache
}

func (ns *nookServer) handleGetLeaderboard(w http.ResponseWriter, r *http.Request) {
	kind, err := models.ParseLeaderboard(r.PathValue("kind"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, err := ns.cache.GetCurrUsers(kind)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (ns *nookServer) handleGetPrevLeaderboard(w http.ResponseWriter, r *http.Request) {
	kind, err := models.ParseLeaderboard(r.PathValue("kind"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users, err := ns.cache.GetPrevUsers(kind)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (ns *nookServer) handleGetUser(w http.ResponseWriter, r *http.Request) {
	kind, err := models.ParseLeaderboard(r.PathValue("kind"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rankParam := r.PathValue("rank")
	rank, err := strconv.Atoi(rankParam)
	if err != nil {
		http.Error(w, fmt.Sprintf(
			"invalid rank: %q conversion to int failed", rankParam,
		), http.StatusBadRequest)
		return
	}

	user, err := ns.cache.GetCurrUser(kind, rank)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (ns *nookServer) handleGetPrevUser(w http.ResponseWriter, r *http.Request) {
	kind, err := models.ParseLeaderboard(r.PathValue("kind"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rankParam := r.PathValue("rank")
	rank, err := strconv.Atoi(rankParam)
	if err != nil {
		http.Error(w, fmt.Sprintf(
			"invalid rank: %q conversion to int failed", rankParam,
		), http.StatusBadRequest)
		return
	}

	user, err := ns.cache.GetPrevUser(kind, rank)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

type ComparisonResponse struct {
	Lower  *models.User `json:"lower"`
	Higher *models.User `json:"higher"`
}

func (ns *nookServer) handleCompareUsers(w http.ResponseWriter, r *http.Request) {
	kind, err := models.ParseLeaderboard(r.PathValue("kind"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rankParams := []string{"rank1", "rank2"}
	var ranks []int

	for _, rankParam := range rankParams {
		param := r.PathValue(rankParam)
		rank, err := strconv.Atoi(param)

		// Check that conversion was successful
		if err != nil {
			http.Error(w, fmt.Sprintf(
				"invalid %s: %q conversion to int failed", rankParam, param,
			), http.StatusBadRequest)
			return
		}

		ranks = append(ranks, rank)
	}

	lowerRank := slices.Max(ranks)
	higherRank := slices.Min(ranks)

	lower, err := ns.cache.GetCurrUser(kind, lowerRank)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	higher, err := ns.cache.GetCurrUser(kind, higherRank)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ComparisonResponse{
		Lower:  lower,
		Higher: higher,
	})
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	lbCache := &cache.LeaderboardCache{}
	interval := cache.ParseRefreshInterval()
	lbCache.LaunchRefresher(ctx, interval)

	ns := &nookServer{cache: lbCache}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /leaderboard/{kind}", ns.handleGetLeaderboard)
	mux.HandleFunc("GET /leaderboard/{kind}/prev", ns.handleGetPrevLeaderboard)
	mux.HandleFunc("GET /leaderboard/{kind}/compare/{rank1}/{rank2}", ns.handleCompareUsers)
	mux.HandleFunc("GET /leaderboard/{kind}/{rank}", ns.handleGetUser)
	mux.HandleFunc("GET /leaderboard/{kind}/prev/{rank}", ns.handleGetPrevUser)

	httpServer := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		slog.Info("nook server listening", "addr", httpServer.Addr, "refresh_interval", interval)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", "err")
		}
	}()

	<-ctx.Done() // block until ctrl-c or SIGTERM
	slog.Info("shutting down nook server")
	httpServer.Shutdown(context.Background())
}
