package v1

import "math/rand"

type Solver struct {
	Game  Game
	Turn  int
	Board Board
	You   Battlesnake
}

// Next returns the next direction that should be applied
// for the given game state.  Optionally, it can return something to yell.
func (s Solver) Next() (Direction, string) {

	// Derive possible moves from given position
	// Takes walls, hazards, own body into consideration
	myPossibleMoves, err := s.You.PossibleMoves(s.Board)
	if err != nil {
		// bleh. Nothing to do.
		return UP, s.negativeResponse()
	}

	// Consider other snakes positions and possible next positions
	otherSnakesPositions := CoordList{}
	for _, snake := range s.Board.Snakes {
		otherSnakesPositions = append(otherSnakesPositions, snake.Body...)

		// Determine possible moves of this snake
		pm, err := snake.PossibleMoves(s.Board)
		if err != nil {
			// snakes next moves is not a threat -- has no valid moves
			continue
		}
		otherSnakesPositions = append(otherSnakesPositions, pm...)
	}

	// Determine if any valid (safe) moves exist
	myPossibleMoves = myPossibleMoves.Eliminate(otherSnakesPositions)

	if len(myPossibleMoves) == 0 {
		// bleh. out of possibilities
		return UP, s.negativeResponse()
	}

	// Have something we can do!
	return myPossibleMoves[rand.Intn(len(myPossibleMoves))].Direction, ""
}

func (s Solver) negativeResponse() string {
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
