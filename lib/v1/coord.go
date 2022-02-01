package v1

import (
	"math"
	"math/rand"
)

// Coord represents a point on the game board, optionally also
// containing the direction that resulted in the Coord.
type Coord struct {
	X         int       `json:"x"`
	Y         int       `json:"y"`
	Direction Direction `json:"-"`
	Score     float64   `json:"-"`
}

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

// WrapForBoard wraps the given coordinate around the edge
// of the given game board so that it remains valid.
// Note this is only valid for the 'wrapped' ruleset.
func (c Coord) WrapForBoard(b Board) Coord {
	if c.X < 0 {
		c.X += b.Width
	} else if c.X > (b.Width - 1) {
		c.X -= b.Width
	}
	if c.Y < 0 {
		c.Y += b.Height
	} else if c.Y > (b.Height - 1) {
		c.Y -= b.Height
	}
	return c
}

// DistanceFrom returns the distance this Coord is from the given Coord
func (c Coord) DistanceFrom(other Coord) float64 {
	return math.Sqrt(math.Pow(float64(other.X-c.X), 2) + math.Pow(float64(other.Y-c.Y), 2))
}

// WithinBounds determines if the Coord is within the boundaries
// of the given Board.
// The game board is represented by a standard 2D grid, oriented with (0,0) in the bottom left.
// The Y-Axis is positive in the up direction, and X-Axis is positive to the right.
// Coordinates begin at zero, such that a board that is 11x11 will have coordinates ranging from [0, 10].
func (c Coord) WithinBounds(b Board) bool {
	return 0 <= c.X && c.X < b.Width && 0 <= c.Y && c.Y < b.Height
}

// CoordList is a list of Coords
type CoordList []Coord

// Contains determines if the given Coord is present in the list
func (cl CoordList) Contains(c Coord) bool {
	if cl == nil || len(cl) == 0 {
		return false
	}
	for _, v := range cl {
		if v.X == c.X && v.Y == c.Y {
			return true
		}
	}
	return false
}

// Eliminate returns a subset CoordList by eliminating any Coord
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

// Rand returns a random coordinate from the list
func (cl CoordList) Rand() Coord {
	return cl[rand.Intn(len(cl))]
}

// Directions returns a slice of Direction referenced by the list
func (cl CoordList) Directions() []Direction {
	ret := []Direction{}
	for _, c := range cl {
		if c.Direction != "" {
			ret = append(ret, c.Direction)
		}
	}
	return ret
}

// Return a subset list made of the first N items
func (cl CoordList) First(n int) CoordList {
	cnt := 0
	resp := CoordList{}
	for _, c := range cl {
		if cnt >= n {
			continue
		}
		resp = append(resp, c)
		cnt += 1
	}
	return resp
}

// AverageDistance returns the average distance the given coordinate is
// from the list by calculating the distance from each point and averaging
// them.
func (cl CoordList) AverageDistance(c Coord) float64 {
	var total float64
	for _, i := range cl {
		total += c.DistanceFrom(i)
	}
	return total / float64(len(cl))
}
