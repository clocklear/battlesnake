package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
)

func router(h *handler) http.Handler {
	r := mux.NewRouter()
	r.Use(nrgorilla.Middleware(h.nr))

	// Wire handlers
	r.HandleFunc("/", h.health).Methods(http.MethodGet)
	r.HandleFunc("/start", h.start).Methods(http.MethodPost)
	r.HandleFunc("/move", h.move).Methods(http.MethodPost)
	r.HandleFunc("/end", h.end).Methods(http.MethodPost)

	return r
}
