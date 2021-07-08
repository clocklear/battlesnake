package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	v1 "github.com/clocklear/battlesnake/lib/v1"
)

type BattlesnakeInfoResponse struct {
	APIVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
}

type GameRequest struct {
	Game  v1.Game        `json:"game"`
	Turn  int            `json:"turn"`
	Board v1.Board       `json:"board"`
	You   v1.Battlesnake `json:"you"`
}

type MoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout,omitempty"`
}

// HandleIndex is called when your Battlesnake is created and refreshed
// by play.battlesnake.com. BattlesnakeInfoResponse contains information about
// your Battlesnake, including what it should look like on the game board.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "clocklear",
		Color:      "#238270",
		Head:       "silly",
		Tail:       "coffee",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

// HandleStart is called at the start of each game your Battlesnake is playing.
// The GameRequest object contains information about the game that's about to start.
// TODO: Use this function to decide how your Battlesnake is going to look on the board.
func HandleStart(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Print("START\n")
}

// HandleMove is called for each turn of each game.
// Valid responses are "up", "down", "left", or "right".
func HandleMove(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Create a solver and use it to determine what we do next
	s := v1.Solver{
		Game:  request.Game,
		Turn:  request.Turn,
		Board: request.Board,
		You:   request.You,
	}

	var resp MoveResponse
	possibleMoves, err := s.Next(v1.SolveOptions{
		UseScoring: true, // enable scoring to optimize next best option
	})
	if err != nil {
		resp.Move = "up"
		resp.Shout = negativeResponse()
	} else {
		resp.Move = string(randDirection(possibleMoves))
	}

	fmt.Printf("GameID: %s - Move: %s\n", request.Game.ID, resp.Move)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Fatal(err)
	}
}

// HandleEnd is called when a game your Battlesnake was playing has ended.
// It's purely for informational purposes, no response required.
func HandleEnd(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Print("END\n")
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/start", HandleStart)
	http.HandleFunc("/move", HandleMove)
	http.HandleFunc("/end", HandleEnd)

	fmt.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func negativeResponse() string {
	r := []string{
		"oh crap",
		"bummer",
		"ouch",
		"whoops",
		"dangit",
		"good game",
		"sayonara",
		"eeeks",
	}
	return r[rand.Intn(len(r))]
}

func randDirection(d []v1.Direction) v1.Direction {
	return d[rand.Intn(len(d))]
}
