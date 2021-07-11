package v1

import "sort"

type Solver struct {
	Game  Game
	Turn  int
	Board Board
	You   Battlesnake
}

type SolveOptions struct {
	UseScoring bool
}

// PossibleMoves returns a list of possible moves that could be taken next
// for the given game state.  An error is raised if something prevents that.
func (s Solver) PossibleMoves(opts SolveOptions) (CoordList, error) {

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
		// My snake can be in this list.  Skip it.
		if snake.ID == s.You.ID {
			continue
		}

		// Gather position of this snakes body pieces
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

	if opts.UseScoring {
		// Score the results
		myPossibleMoves = s.score(myPossibleMoves).First(2)
	}

	// Return the results
	return myPossibleMoves, nil
}

func (s Solver) score(moves CoordList) CoordList {
	// Given the list of possible moves, 'score' each one, sort the list
	// based on score, and return
	scored := CoordList{}
	for _, m := range moves {
		// Score by avoiding self
		// Find avg distance to first 8 body points
		avgDistance := s.You.Body.First(8).AverageDistance(m)
		m.Score = avgDistance

		// Amend score by considering food
		// If our health is above 70 and this move overlaps food, avoid it
		isFood := s.Board.Food.Contains(m)
		if s.You.Health >= 70 && isFood {
			m.Score -= 5
		}
		// If our health is below 30 and this move overlaps food, bump it in priority
		if s.You.Health <= 30 && isFood {
			m.Score += 5
		}
		scored = append(scored, m)
	}

	// Sort the result
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score < scored[j].Score
	})

	return scored
}

func (s Solver) PickMove(possibleMoves CoordList) (Direction, error) {
	switch len(possibleMoves) {
	case 0:
		return randDirection(allDirections), ErrNoPossibleMove
	case 1:
		return possibleMoves[0].Direction, nil
	default:
		// If the first option here is significantly stronger than the others, use it
		if possibleMoves[0].Score-possibleMoves[1].Score >= 4 {
			return possibleMoves[0].Direction, nil
		}
		// Otherwise pick randomly from first two items
		return randDirection(possibleMoves.First(2).Directions()), nil
	}
}
