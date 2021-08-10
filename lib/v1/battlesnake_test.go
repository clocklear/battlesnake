package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPossibleMoves(t *testing.T) {
	testCases := []struct {
		desc          string
		s             Battlesnake
		b             Board
		expected      CoordList
		expectedError error
	}{
		{
			desc: "success, simple",
			s: Battlesnake{
				Head: Coord{
					X: 5,
					Y: 5,
				},
				Body: CoordList{
					{
						X: 5,
						Y: 5,
					},
				},
			},
			b: Board{
				Height: 11,
				Width:  11,
			},
			expected: CoordList{
				{
					X:         5,
					Y:         6,
					Direction: UP,
				},
				{
					X:         5,
					Y:         4,
					Direction: DOWN,
				},
				{
					X:         4,
					Y:         5,
					Direction: LEFT,
				},
				{
					X:         6,
					Y:         5,
					Direction: RIGHT,
				},
			},
		},
		{
			desc: "success, cannot go up because body",
			s: Battlesnake{
				Head: Coord{
					X: 5,
					Y: 5,
				},
				Body: CoordList{
					{
						X: 5,
						Y: 5,
					},
					{
						X: 5,
						Y: 6,
					},
				},
			},
			b: Board{
				Height: 11,
				Width:  11,
			},
			expected: CoordList{
				{
					X:         5,
					Y:         4,
					Direction: DOWN,
				},
				{
					X:         4,
					Y:         5,
					Direction: LEFT,
				},
				{
					X:         6,
					Y:         5,
					Direction: RIGHT,
				},
			},
		},
		{
			desc: "success, must go right because walls and body",
			s: Battlesnake{
				Head: Coord{
					X: 0,
					Y: 0,
				},
				Body: CoordList{
					{
						X: 0,
						Y: 0,
					},
					{
						X: 0,
						Y: 1,
					},
				},
			},
			b: Board{
				Height: 11,
				Width:  11,
			},
			expected: CoordList{
				{
					X:         1,
					Y:         0,
					Direction: RIGHT,
				},
			},
		},
		{
			desc: "failure, no moves due to walls and body (bottom left)",
			s: Battlesnake{
				Head: Coord{
					X: 0,
					Y: 0,
				},
				Body: CoordList{
					{
						X: 0,
						Y: 0,
					},
					{
						X: 0,
						Y: 1,
					},
					{
						X: 1,
						Y: 1,
					},
					{
						X: 1,
						Y: 0,
					},
				},
			},
			b: Board{
				Height: 11,
				Width:  11,
				Hazards: CoordList{
					{
						X: 1,
						Y: 0,
					},
				},
			},
			expectedError: ErrNoPossibleMove,
		},

		{
			desc: "success, must go left because walls and body",
			s: Battlesnake{
				Head: Coord{
					X: 10,
					Y: 10,
				},
				Body: CoordList{
					{
						X: 10,
						Y: 10,
					},
					{
						X: 10,
						Y: 9,
					},
				},
			},
			b: Board{
				Height: 11,
				Width:  11,
			},
			expected: CoordList{
				{
					X:         9,
					Y:         10,
					Direction: LEFT,
				},
			},
		},
		{
			desc: "failure, no moves due to walls and body",
			s: Battlesnake{
				Head: Coord{
					X: 10,
					Y: 10,
				},
				Body: CoordList{
					{
						X: 10,
						Y: 10,
					},
					{
						X: 10,
						Y: 9,
					},
					{
						X: 9,
						Y: 9,
					},
					{
						X: 9,
						Y: 10,
					},
				},
			},
			b: Board{
				Height: 11,
				Width:  11,
				Hazards: CoordList{
					{
						X: 9,
						Y: 10,
					},
				},
			},
			expectedError: ErrNoPossibleMove,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual, err := tC.s.PossibleMoves(tC.b)
			if tC.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tC.expectedError)
				return
			}
			assert.ElementsMatch(t, tC.expected, actual)
		})
	}
}
