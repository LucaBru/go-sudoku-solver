package sequential

import (
	"sudoku/sudoku"
)

func Solver(puzzle sudoku.Puzzle, row, clm int) bool {
	if clm > 8 {
		row++
		clm = 0
	}
	if row == 9 {
		return true
	}
	if len(puzzle[row][clm]) == 1 {
		return Solver(puzzle, row, clm+1)
	}

	candidates := map[int]struct{}{}
	for k := range puzzle[row][clm] {
		candidates[k] = struct{}{}
	}

	for d := range candidates {
		if puzzle.Valid(d, row, clm) {
			puzzle[row][clm] = map[int]struct{}{d: {}}
			if Solver(puzzle, row, clm+1) {
				return true
			}
		}
	}
	puzzle[row][clm] = candidates
	return false
}
