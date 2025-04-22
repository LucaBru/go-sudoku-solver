package sequential

import (
	"sudoku/sudoku"
)

func Solver(puzzle sudoku.Puzzle, pos sudoku.Pos) sudoku.Puzzle {
	if pos.Row == 9 {
		return puzzle
	}
	if ok, _ := puzzle[pos.Row][pos.Clm].IsSingleton(); ok {
		return Solver(puzzle, *pos.Next())
	}
	_, values := puzzle[pos.Row][pos.Clm].IsSingleton()
	for _, v := range values {
		if puzzle.Valid(v, pos) {
			puzzle[pos.Row][pos.Clm].SetSingleton(v)
			if sol := Solver(puzzle, *pos.Next()); sol != nil {
				return sol
			}
		}
	}
	puzzle[pos.Row][pos.Clm] = sudoku.NewDigits()
	return nil 
}
