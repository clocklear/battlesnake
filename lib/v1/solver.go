package v1

type Solver struct {
	Game  Game
	Turn  int
	Board Board
	You   Battlesnake
}

// Next returns a list of possible directions that could be taken next
// for the given game state.  An error is raised if something prevents that.
func (s Solver) Next() ([]Direction, error) {

	// Derive possible moves from given position
	// Takes walls, hazards, own body into consideration
	myPossibleMoves, err := s.You.PossibleMoves(s.Board)
	if err != nil {
		// bleh. Nothing to do.
		return nil, err
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
		return nil, ErrNoPossibleMove
	}

	// Have things we can do!
	return myPossibleMoves.Directions(), nil
}
