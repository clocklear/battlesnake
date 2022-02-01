package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordProject(t *testing.T) {
	testCases := []struct {
		desc     string
		starting Coord
		dir      Direction
		expected Coord
	}{
		{
			desc: "UP",
			starting: Coord{
				X: 5,
				Y: 5,
			},
			dir: UP,
			expected: Coord{
				X:         5,
				Y:         6,
				Direction: UP,
			},
		},
		{
			desc: "DOWN",
			starting: Coord{
				X: 5,
				Y: 5,
			},
			dir: DOWN,
			expected: Coord{
				X:         5,
				Y:         4,
				Direction: DOWN,
			},
		},
		{
			desc: "LEFT",
			starting: Coord{
				X: 5,
				Y: 5,
			},
			dir: LEFT,
			expected: Coord{
				X:         4,
				Y:         5,
				Direction: LEFT,
			},
		},
		{
			desc: "RIGHT",
			starting: Coord{
				X: 5,
				Y: 5,
			},
			dir: RIGHT,
			expected: Coord{
				X:         6,
				Y:         5,
				Direction: RIGHT,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.starting.Project(tC.dir)
			assert.Equal(t, tC.expected, actual)
		})
	}
}

func TestCoordWithinBounds(t *testing.T) {
	testCases := []struct {
		desc     string
		c        []Coord
		b        Board
		expected bool
	}{
		{
			desc: "success",
			c: []Coord{
				{
					X: 1,
					Y: 1,
				},
				{
					X: 6,
					Y: 2,
				},
				{
					X: 0,
					Y: 0,
				},
				{
					X: 9,
					Y: 9,
				},
			},
			b: Board{
				Height: 10,
				Width:  10,
			},
			expected: true,
		},
		{
			desc: "failure",
			c: []Coord{
				{
					X: -1,
					Y: 5,
				},
				{
					X: 12,
					Y: 12,
				},
				{
					X: -5,
					Y: -1,
				},
			},
			b: Board{
				Height: 10,
				Width:  10,
			},
			expected: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			for _, c := range tC.c {
				actual := c.WithinBounds(tC.b)
				assert.Equal(t, tC.expected, actual)
			}
		})
	}
}

func TestCoordListContains(t *testing.T) {
	testCases := []struct {
		desc     string
		list     CoordList
		coords   []Coord
		expected bool
	}{
		{
			desc: "success, list contains items",
			list: CoordList{
				{
					X: 1,
					Y: 1,
				},
				{
					X: 2,
					Y: 2,
				},
				{
					X: 3,
					Y: 3,
				},
			},
			coords: []Coord{
				{
					X: 1,
					Y: 1,
				},
				{
					X: 2,
					Y: 2,
				},
				{
					X: 3,
					Y: 3,
				},
			},
			expected: true,
		},
		{
			desc: "failure, list missing items",
			list: CoordList{
				{
					X: 1,
					Y: 1,
				},
				{
					X: 2,
					Y: 2,
				},
				{
					X: 3,
					Y: 3,
				},
			},
			coords: []Coord{
				{
					X: 1,
					Y: 2,
				},
				{
					X: 2,
					Y: 3,
				},
				{
					X: 3,
					Y: 4,
				},
			},
			expected: false,
		},
		{
			desc: "failure, empty list",
			list: CoordList{},
			coords: []Coord{
				{
					X: 1,
					Y: 2,
				},
				{
					X: 2,
					Y: 3,
				},
				{
					X: 3,
					Y: 4,
				},
			},
			expected: false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			for _, c := range tC.coords {
				actual := tC.list.Contains(c)
				assert.Equal(t, tC.expected, actual)
			}
		})
	}
}

func TestCoordListEliminate(t *testing.T) {
	testCases := []struct {
		desc      string
		start     CoordList
		eliminate CoordList
		expected  CoordList
	}{
		{
			desc: "overlap, eliminated item",
			start: CoordList{
				{
					X: 1,
					Y: 1,
				},
				{
					X: 2,
					Y: 2,
				},
			},
			eliminate: CoordList{
				{
					X: 1,
					Y: 1,
				},
				{
					X: 3,
					Y: 3,
				},
			},
			expected: CoordList{
				{
					X: 2,
					Y: 2,
				},
			},
		},
		{
			desc: "no overlap",
			start: CoordList{
				{
					X: 1,
					Y: 1,
				},
				{
					X: 2,
					Y: 2,
				},
			},
			eliminate: CoordList{
				{
					X: 3,
					Y: 3,
				},
			},
			expected: CoordList{
				{
					X: 1,
					Y: 1,
				},
				{
					X: 2,
					Y: 2,
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.start.Eliminate(tC.eliminate)
			assert.NotSame(t, tC.start, actual)
			assert.ElementsMatch(t, tC.expected, actual)
		})
	}
}

func TestCoordListWrapForBoard(t *testing.T) {
	testCases := []struct {
		desc     string
		in       Coord
		b        Board
		expected Coord
	}{
		{
			desc: "standard coord, at origin",
			in: Coord{
				X: 0,
				Y: 0,
			},
			b: Board{
				Height: 11,
				Width:  11,
			},
			expected: Coord{
				X: 0,
				Y: 0,
			},
		},
		{
			desc: "standard coord, at extents",
			in: Coord{
				X: 10,
				Y: 10,
			},
			b: Board{
				Height: 11,
				Width:  11,
			},
			expected: Coord{
				X: 10,
				Y: 10,
			},
		},
		{
			desc: "extended coord, high beyond X",
			in: Coord{
				X: 11,
				Y: 10,
			},
			b: Board{
				Height: 11,
				Width:  11,
			},
			expected: Coord{
				X: 0,
				Y: 10,
			},
		},
		{
			desc: "extended coord, high beyond Y",
			in: Coord{
				X: 10,
				Y: 11,
			},
			b: Board{
				Height: 11,
				Width:  11,
			},
			expected: Coord{
				X: 10,
				Y: 0,
			},
		},
		{
			desc: "extended coord, below X",
			in: Coord{
				X: -1,
				Y: 5,
			},
			b: Board{
				Height: 11,
				Width:  11,
			},
			expected: Coord{
				X: 10,
				Y: 5,
			},
		},
		{
			desc: "extended coord, below Y",
			in: Coord{
				X: 5,
				Y: -1,
			},
			b: Board{
				Height: 11,
				Width:  11,
			},
			expected: Coord{
				X: 5,
				Y: 10,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.in.WrapForBoard(tC.b)
			assert.Equal(t, tC.expected, actual)
		})
	}
}
