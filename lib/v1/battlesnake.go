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

const (
	HazardDamagePerTurn = 15
	MaximumSnakeHealth  = 100
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
func (bs Battlesnake) PossibleMoves(b Board, g Game) (CoordList, error) {
	cl := CoordList{}
	// for each direction..
	for _, d := range allDirections {
		// ..project the head that way
		c := bs.Head.Project(d)

		// Projection is... different for wrapped games
		if g.Ruleset.Name == RulesetWrapped {
			// We need to wrap this coord in a wrapped game
			c = c.WrapForBoard(b)
		} else {
			// Does it fit on the board?
			if !c.WithinBounds(b) {
				continue
			}
		}

		// Does it overlap with our body?
		if bs.Body.Contains(c) {
			continue
		}

		// Looks like a valid move
		// Might be a hazard, but it's still valid.
		cl = append(cl, c)
	}
	if len(cl) == 0 {
		return cl, ErrNoPossibleMove
	}
	return cl, nil
}

// Project moves the Battlesnake to the given coordinate on the given board
func (bs Battlesnake) Project(loc Coord, board Board) Battlesnake {
	bs.Head = loc
	// Prepend head to body
	bs.Body = append([]Coord{loc}, bs.Body...)
	// Decrement health
	bs.Health -= 1
	willGrow := board.Food.Contains(loc)
	if willGrow {
		bs.Health = MaximumSnakeHealth
	} else {
		// Drop last elem
		bs.Body = bs.Body[:len(bs.Body)-1]
	}
	if board.Hazards.Contains(loc) {
		bs.Health -= HazardDamagePerTurn
	}
	return bs
}

// IsValid determines if the snake is 'valid' on the given board.
// Valid snakes:
// * have possible moves
// * have non-zero health
func (bs Battlesnake) IsValid(b Board, g Game) bool {
	_, err := bs.PossibleMoves(b, g)
	return err == nil && bs.Health > 0
}
