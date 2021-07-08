package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tstGame Game = Game{
	ID:      "tst",
	Timeout: 500,
}

var simpleEmptyBoard Board = Board{
	Height: 11,
	Width:  11,
}

var opposingSnakeBoard Board = Board{
	Height: 11,
	Width:  11,
	Snakes: []Battlesnake{
		{
			Head: Coord{
				X: 1,
				Y: 1,
			},
			Body: CoordList{
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
	},
}

func TestNext(t *testing.T) {
	testCases := []struct {
		desc               string
		game               Game
		turn               int
		board              Board
		you                Battlesnake
		possibleDirections []Direction
		expectedError      error
	}{
		{
			desc:  "simple",
			game:  tstGame,
			board: simpleEmptyBoard,
			you: Battlesnake{
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
			possibleDirections: []Direction{
				UP,
				DOWN,
				LEFT,
				RIGHT,
			},
		},
		{
			desc:  "body limits options",
			game:  tstGame,
			board: simpleEmptyBoard,
			you: Battlesnake{
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
						X: 6,
						Y: 5,
					},
				},
			},
			possibleDirections: []Direction{
				UP,
				DOWN,
				LEFT,
			},
		},
		{
			desc:  "walls limit options",
			game:  tstGame,
			board: simpleEmptyBoard,
			you: Battlesnake{
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
						X: 1,
						Y: 0,
					},
				},
			},
			possibleDirections: []Direction{
				UP,
			},
		},
		{
			desc: "no moves available",
			game: tstGame,
			board: Board{
				Width:  5,
				Height: 5,
				Hazards: CoordList{
					{
						X: 0,
						Y: 1,
					},
				},
			},
			you: Battlesnake{
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
						X: 1,
						Y: 0,
					},
				},
			},
			expectedError: ErrNoPossibleMove,
		},
		{
			desc:  "snake limits options",
			game:  tstGame,
			board: opposingSnakeBoard,
			you: Battlesnake{
				Head: Coord{
					X: 0,
					Y: 1,
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
			possibleDirections: []Direction{
				UP,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			s := Solver{
				Game:  tC.game,
				Turn:  tC.turn,
				Board: tC.board,
				You:   tC.you,
			}
			possibleMoves, err := s.Next()
			if tC.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tC.expectedError)
				return
			}
			assert.ElementsMatch(t, possibleMoves, tC.possibleDirections)
		})
	}
}
