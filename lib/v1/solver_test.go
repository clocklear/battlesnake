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
			ID: "not-you",
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

func TestSolverPossibleMoves(t *testing.T) {
	testCases := []struct {
		desc               string
		game               Game
		turn               int
		board              Board
		you                Battlesnake
		possibleDirections []Direction
		opts               SolveOptions
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
			opts: SolveOptions{},
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
			opts: SolveOptions{},
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
			opts: SolveOptions{},
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
			opts:          SolveOptions{},
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
			opts: SolveOptions{},
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
			possibleMoves, err := s.PossibleMoves(tC.opts)
			if tC.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tC.expectedError)
				return
			}
			assert.ElementsMatch(t, possibleMoves.Directions(), tC.possibleDirections)
		})
	}
}

func TestSolverPickMove(t *testing.T) {
	testCases := []struct {
		desc          string
		possibleMoves CoordList
		expected      []Direction
		expectedError error
	}{
		{
			desc: "multiple possibilities, but clear winner",
			possibleMoves: CoordList{
				{
					X:         5,
					Y:         6,
					Direction: UP,
					Score:     7,
				},
				{
					X:         4,
					Y:         5,
					Direction: LEFT,
					Score:     3,
				},
				{
					X:         6,
					Y:         5,
					Direction: RIGHT,
					Score:     2.5,
				},
			},
			expected: []Direction{UP},
		},
		{
			desc: "one possibility",
			possibleMoves: CoordList{
				{
					X:         6,
					Y:         5,
					Direction: RIGHT,
					Score:     2.5,
				},
			},
			expected: []Direction{RIGHT},
		},
		{
			desc:          "no moves",
			possibleMoves: CoordList{},
			expectedError: ErrNoPossibleMove,
		},
		{
			desc: "multiple possibilities, no clear winner",
			possibleMoves: CoordList{
				{
					X:         5,
					Y:         6,
					Direction: UP,
					Score:     1,
				},
				{
					X:         4,
					Y:         5,
					Direction: LEFT,
					Score:     3,
				},
				{
					X:         6,
					Y:         5,
					Direction: RIGHT,
					Score:     2.5,
				},
			},
			expected: []Direction{LEFT, RIGHT},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			s := Solver{}
			actual, err := s.PickMove(tC.possibleMoves)
			if tC.expectedError != nil {
				assert.ErrorIs(t, err, tC.expectedError)
				return
			}
			assert.Contains(t, tC.expected, actual)
		})
	}
}

func TestSolverScore(t *testing.T) {
	testCases := []struct {
		desc     string
		game     Game
		turn     int
		board    Board
		you      Battlesnake
		moves    CoordList
		expected CoordList
	}{
		{
			desc:  "body length 1, no food",
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
			moves: CoordList{
				{
					X:         6,
					Y:         5,
					Direction: RIGHT,
				},
				{
					X:         4,
					Y:         5,
					Direction: LEFT,
				},
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
			},
			expected: CoordList{
				{
					X:         6,
					Y:         5,
					Direction: RIGHT,
					Score:     1,
				},
				{
					X:         4,
					Y:         5,
					Direction: LEFT,
					Score:     1,
				},
				{
					X:         5,
					Y:         6,
					Direction: UP,
					Score:     1,
				},
				{
					X:         5,
					Y:         4,
					Direction: DOWN,
					Score:     1,
				},
			},
		},
		{
			desc:  "body length 5, no food",
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
						X: 5,
						Y: 4,
					},
					{
						X: 5,
						Y: 3,
					},
					{
						X: 5,
						Y: 2,
					},
					{
						X: 5,
						Y: 1,
					},
				},
			},
			moves: CoordList{
				{
					X:         6,
					Y:         5,
					Direction: RIGHT,
				},
				{
					X:         4,
					Y:         5,
					Direction: LEFT,
				},
				{
					X:         5,
					Y:         6,
					Direction: UP,
				},
			},
			expected: CoordList{
				{
					X:         5,
					Y:         6,
					Direction: UP,
					Score:     3,
				},
				{
					X:         6,
					Y:         5,
					Direction: RIGHT,
					Score:     2.387132965131785,
				},
				{
					X:         4,
					Y:         5,
					Direction: LEFT,
					Score:     2.387132965131785,
				},
			},
		},
		{
			desc: "body length 5, food nearby but not hungry",
			game: tstGame,
			board: Board{
				Height: 11,
				Width:  11,
				Food: CoordList{
					{
						X: 5,
						Y: 6,
					},
				},
			},
			you: Battlesnake{
				Health: 100,
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
						Y: 4,
					},
					{
						X: 5,
						Y: 3,
					},
					{
						X: 5,
						Y: 2,
					},
					{
						X: 5,
						Y: 1,
					},
				},
			},
			moves: CoordList{
				{
					X:         6,
					Y:         5,
					Direction: RIGHT,
				},
				{
					X:         4,
					Y:         5,
					Direction: LEFT,
				},
				{
					X:         5,
					Y:         6,
					Direction: UP,
				},
			},
			expected: CoordList{
				{
					X:         6,
					Y:         5,
					Direction: RIGHT,
					Score:     2.387132965131785,
				},
				{
					X:         4,
					Y:         5,
					Direction: LEFT,
					Score:     2.387132965131785,
				},
				{
					X:         5,
					Y:         6,
					Direction: UP,
					Score:     -2,
				},
			},
		},
		{
			desc: "body length 5, food nearby and starving",
			game: tstGame,
			board: Board{
				Height: 11,
				Width:  11,
				Food: CoordList{
					{
						X: 6,
						Y: 5,
					},
				},
			},
			you: Battlesnake{
				Health: 25,
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
						Y: 4,
					},
					{
						X: 5,
						Y: 3,
					},
					{
						X: 5,
						Y: 2,
					},
					{
						X: 5,
						Y: 1,
					},
				},
			},
			moves: CoordList{
				{
					X:         6,
					Y:         5,
					Direction: RIGHT,
				},
				{
					X:         4,
					Y:         5,
					Direction: LEFT,
				},
				{
					X:         5,
					Y:         6,
					Direction: UP,
				},
			},
			expected: CoordList{
				{
					X:         6,
					Y:         5,
					Direction: RIGHT,
					Score:     7.387132965131785,
				},
				{
					X:         5,
					Y:         6,
					Direction: UP,
					Score:     3,
				},
				{
					X:         4,
					Y:         5,
					Direction: LEFT,
					Score:     2.387132965131785,
				},
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
			scoredMoves := s.score(tC.moves)
			assert.Equal(t, tC.expected, scoredMoves) // order is important
		})
	}
}
