package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/clocklear/battlesnake/lib/gamerecorder"
	v1 "github.com/clocklear/battlesnake/lib/v1"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type handler struct {
	rec gamerecorder.GameRecorder
	l   logger
	nr  *newrelic.Application
	so  v1.SolveOptions // TODO: this probably shouldn't live here
}

type BattlesnakeInfoResponse struct {
	Author string `json:"author"`
	Color  string `json:"color"`
	Head   string `json:"head"`
	Tail   string `json:"tail"`
}

func (h *handler) health(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())
	response := BattlesnakeInfoResponse{
		Author: "clocklear",
		Color:  "#238270",
		Head:   "silly",
		Tail:   "coffee",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		txn.NoticeError(err)
		h.l.Error("failed to encode health response", "err", err.Error())
	}
}

func (h *handler) start(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())
	request := v1.GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.l.Error("bad start request", "err", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.rec.Start(context.Background(), request)
	if err != nil {
		txn.NoticeError(err)
		h.l.Error("failed to record game start", "err", err.Error(), "gameId", request.Game.ID)
	}

	// Nothing to respond with here
	// h.l.Info("starting game", "gameId", request.Game.ID)
}

type moveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout,omitempty"`
}

func (h *handler) CreateMoveHandlerWithSolveOpts(opts v1.SolveOptions) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		txn := newrelic.FromContext(r.Context())
		request := v1.GameRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			h.l.Error("bad move request", "err", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Create a solver and use it to determine what we do next
		s := v1.CreateSolver(request).WithLogger(h.l.base)

		var resp moveResponse
		possibleMoves, err := s.PossibleMoves(opts)
		var move string
		if err != nil {
			resp.Move = "up"
			move = "invalid"
		} else {
			d, _ := s.PickMove(possibleMoves, opts)
			resp.Move = string(d)
			move = resp.Move
		}

		// Record this move
		err = h.rec.Move(context.Background(), request, move)
		if err != nil {
			txn.NoticeError(err)
			h.l.Error("failed to record game move", "game", request.Game.ID, "turn", request.Turn, "move", resp.Move, "err", err.Error())
		}
		// h.l.Info("responding with move", "game", request.Game.ID, "move", resp.Move)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			txn.NoticeError(err)
			h.l.Error("failed to encode move response", "err", err.Error())
		}
	}
}

func (h *handler) end(w http.ResponseWriter, r *http.Request) {
	txn := newrelic.FromContext(r.Context())
	request := v1.GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		h.l.Error("bad end game request", "err", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.rec.End(context.Background(), request)
	if err != nil {
		txn.NoticeError(err)
		h.l.Error("failed recording game end", "err", err.Error(), "game", request.Game.ID)
	}

	// Nothing to respond with here
	// h.l.Info("ending game", "gameId", request.Game.ID)
}
