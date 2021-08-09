package v1

import "fmt"

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
func (bs Battlesnake) PossibleMoves(b Board, opts SolveOptions) (CoordList, error) {
	cl := CoordList{}
	// for each direction..
	for _, d := range allDirections {
		// ..project the head that way
		c := bs.Head.Project(d)

		// Give the move a base score of 50
		c.Score = 50

		// Does it fit on the board?
		if !c.WithinBounds(b) {
			continue
		}
		// Does it overlap with our body?
		if bs.Body.Contains(c) {
			continue
		}
		// Does it overlap with food?  Reward!
		if b.Food != nil && len(b.Food) > 0 {
			if b.Food.Contains(c) {
				// Improve the score of food moves
				c.Score += 20
			}
		}

		// Does it overlap with the board hazards?
		if b.Hazards != nil && len(b.Hazards) > 0 {
			if b.Hazards.Contains(c) {
				// Reduce the score of hazard moves
				c.Score -= 40
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

// Project moves the Battlesnake to the given coordinate
func (bs Battlesnake) Project(loc Coord, willGrow bool) Battlesnake {
	bs.Head = loc
	// Prepend head to body
	bs.Body = append([]Coord{loc}, bs.Body...)
	if !willGrow {
		// Drop last elem
		bs.Body = bs.Body[:len(bs.Body)-1]
	}
	return bs
}

// IsValid determines if the snake is 'valid' on the given board.
// Valid snakes:
// * have possible moves
// * have non-zero health
func (bs Battlesnake) IsValid(b Board) bool {
	_, err := bs.PossibleMoves(b, SolveOptions{})
	return err == nil && bs.Health > 0
}
