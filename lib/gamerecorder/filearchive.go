package gamerecorder

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	v1 "github.com/clocklear/battlesnake/lib/v1"
)

// FileArchive is a gamerecorder implementation that writes games to gzipped archives
// in the given path
type FileArchive struct {
	basePath          string
	games             map[string]game
	maxAgeBeforePrune time.Duration
	pruneInterval     time.Duration
	quit              chan int
	mu                sync.RWMutex
}

type decision struct {
	BoardState v1.BoardState `json:"state"`
	Decision   string        `json:"decision"`
}

type game struct {
	Game       v1.Game    `json:"game"`
	Decisions  []decision `json:"states"`
	Started    time.Time  `json:"startedAt"`
	Ended      time.Time  `json:"endedAt"`
	Won        bool       `json:"won"`
	expiration int64
}

func NewFileArchive(basePath string, pruneInterval time.Duration, maxAgeBeforePrune time.Duration) GameRecorder {
	fa := FileArchive{
		games:             make(map[string]game),
		basePath:          basePath,
		pruneInterval:     pruneInterval,
		maxAgeBeforePrune: maxAgeBeforePrune,
		quit:              make(chan int, 1),
	}

	// start prune loop
	go fa.pruneloop()

	return &fa
}

func (r *FileArchive) pruneloop() {
	tick := time.NewTicker(r.pruneInterval)

	for {
		select {
		case <-tick.C:
			r.prune()
		case <-r.quit:
			tick.Stop()
			return
		}
	}
}

func (r *FileArchive) prune() {
	gameIds := []string{}
	r.mu.RLock()
	for gameId, g := range r.games {
		if g.expiration <= time.Now().UnixNano() {
			gameIds = append(gameIds, gameId)
		}
	}
	r.mu.RUnlock()
	if len(gameIds) == 0 {
		return
	}
	for _, g := range gameIds {
		// purge does its own locking
		r.purge(g)
	}
}

func gameKey(req v1.GameRequest) string {
	return fmt.Sprintf("%v:%v", req.Game.ID, req.You.ID)
}

func (r *FileArchive) Start(ctx context.Context, req v1.GameRequest) error {
	// Start a new game
	r.mu.Lock()
	r.games[gameKey(req)] = game{
		Game:       req.Game,
		Decisions:  []decision{},
		Started:    time.Now(),
		expiration: time.Now().Add(r.maxAgeBeforePrune).UnixNano(),
	}
	r.mu.Unlock()
	return nil
}

func (r *FileArchive) Move(ctx context.Context, req v1.GameRequest, move string) error {
	key := gameKey(req)
	r.mu.RLock()
	g, validGame := r.games[key]
	r.mu.RUnlock()
	if !validGame {
		// C. Locklear -- sometimes we don't get a start request and move is invoked immediately
		// Just start a game
		err := r.Start(ctx, req)
		if err != nil {
			return err
		}
		r.mu.RLock()
		g = r.games[key]
		r.mu.RUnlock()
	}
	g.Decisions = append(g.Decisions, decision{
		BoardState: req.ToBoardState(),
		Decision:   move,
	})
	r.mu.Lock()
	r.games[key] = g
	r.mu.Unlock()
	return nil
}

func (r *FileArchive) End(ctx context.Context, req v1.GameRequest) error {
	r.mu.RLock()
	g, validGame := r.games[gameKey(req)]
	r.mu.RUnlock()
	if !validGame {
		return fmt.Errorf("invalid game")
	}
	g.Ended = time.Now()
	g.Decisions = append(g.Decisions, decision{
		BoardState: req.ToBoardState(),
		Decision:   "end",
	})
	g.Won = didWin(g, req)

	// Setup cleanup
	defer r.purge(req.Game.ID)

	// Render game to json
	jsonGame, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}

	// Open a file for writing.
	outputFile := path.Join(r.basePath, fmt.Sprintf("%v_game=%v_type=%v_snake=%v.json.gz", g.Ended.Format("20060102T150405Z"), g.Game.ID, req.Game.Ruleset.Name, req.You.Name))
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	// Create gzip writer.
	w := gzip.NewWriter(f)

	// Write bytes in compressed form to the file.
	_, err = w.Write(jsonGame)
	if err != nil {
		return err
	}

	// Close the file.
	return w.Close()
}

func (r *FileArchive) Shutdown() error {
	r.quit <- 1
	return nil
}

func (r *FileArchive) purge(gameId string) {
	r.mu.Lock()
	delete(r.games, gameId)
	r.mu.Unlock()
}

func didWin(g game, endState v1.GameRequest) bool {
	return endState.You.IsValid(endState.Board, endState.Game) && hasNoInvalidDecision(g)
}

func hasNoInvalidDecision(g game) bool {
	hasInvalidDecision := false
	for _, d := range g.Decisions {
		if d.Decision == "invalid" {
			hasInvalidDecision = true
		}
	}
	return hasInvalidDecision
}
