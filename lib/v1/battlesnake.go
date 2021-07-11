package v1

import (
	"fmt"
	"math/rand"
)

// ErrNoPossibleMove indicates that no legal move is possible from the current position
var ErrNoPossibleMove = fmt.Errorf("no possible moves")

type Direction string

const (
	UP    Direction = "up"
	DOWN  Direction = "down"
	LEFT  Direction = "left"
	RIGHT Direction = "right"
)

var allDirections []Direction = []Direction{
	UP, DOWN, LEFT, RIGHT,
}

func randDirection(d []Direction) Direction {
	return d[rand.Intn(len(d))]
}

type Battlesnake struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Health int32     `json:"health"`
	Body   CoordList `json:"body"`
	Head   Coord     `json:"head"`
	Length int32     `json:"length"`
	Shout  string    `json:"shout"`
}

// PossibleMoves returns the list of possible coords
// the Battlesnake could take based on its current
// position and the provided board.  It takes the board
// bounds and hazards into consideration.  An error
// will the thrown if no moves are possible.
func (bs Battlesnake) PossibleMoves(b Board) (CoordList, error) {
	cl := CoordList{}
	// for each direction..
	for _, d := range allDirections {
		// ..project the head that way
		c := bs.Head.Project(d)
		// Does it fit on the board?
		if !c.WithinBounds(b) {
			continue
		}
		// Does it overlap with our body?
		if bs.Body.Contains(c) {
			continue
		}
		// Does it overlap with the board hazards?
		if b.Hazards != nil && len(b.Hazards) > 0 {
			if b.Hazards.Contains(c) {
				continue
			}
		}
		// Looks like a valid move
		cl = append(cl, c)
	}
	if len(cl) == 0 {
		return cl, ErrNoPossibleMove
	}
	return cl, nil
}
