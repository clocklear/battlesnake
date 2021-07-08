package v1

import "fmt"

type Game struct {
	ID      string `json:"id"`
	Timeout int32  `json:"timeout"`
}

type Coord struct {
	X         int       `json:"x"`
	Y         int       `json:"y"`
	Direction Direction `json:"-"`
}

// NoPossibleMove indicates that no legal move is possible from the current position
var NoPossibleMove = fmt.Errorf("no possible moves")

// Project shifts the Coord in the given direction
// and returns the result.  The resulting position
// is not guaranteed to be valid.
func (c Coord) Project(d Direction) Coord {
	switch d {
	case UP:
		c = Coord{
			X:         c.X,
			Y:         c.Y + 1,
			Direction: UP,
		}
	case DOWN:
		c = Coord{
			X:         c.X,
			Y:         c.Y - 1,
			Direction: DOWN,
		}
	case LEFT:
		c = Coord{
			X:         c.X - 1,
			Y:         c.Y,
			Direction: LEFT,
		}
	case RIGHT:
		c = Coord{
			X:         c.X + 1,
			Y:         c.Y,
			Direction: RIGHT,
		}
	}
	return c
}

// WithinBounds determines if the Coord is within the boundaries
// of the given Board.
func (c Coord) WithinBounds(b Board) bool {
	return 0 <= c.X && c.X < b.Width && 0 <= c.Y && c.Y < b.Height
}

// CoordList is a list of Coords
type CoordList []Coord

// Contains determines if the given Coord is present in the list
func (cl CoordList) Contains(c Coord) bool {
	for _, v := range cl {
		if v.X == c.X && v.Y == c.Y {
			return true
		}
	}
	return false
}

// Eliminate returns a subset CoordList by eliminating matches from
// present in the candidates list
func (cl CoordList) Eliminate(candidates CoordList) CoordList {
	r := CoordList{}
	for _, c := range cl {
		if candidates.Contains(c) {
			continue
		}
		r = append(r, c)
	}
	return r
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
		return cl, NoPossibleMove
	}
	return cl, nil
}

type Board struct {
	Height  int           `json:"height"`
	Width   int           `json:"width"`
	Food    CoordList     `json:"food"`
	Hazards CoordList     `json:"hazards"`
	Snakes  []Battlesnake `json:"snakes"`
}

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
