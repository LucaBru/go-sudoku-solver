package solver

import "sudoku/utils"

func Sequential(puzzle *utils.Puzzle, row, clm int) bool {
	if clm > 8 {
		row++
		clm = 0
	}
	if row == 9 {
		return true
	}

	if len(puzzle[row][clm]) == 1 {
		return Sequential(puzzle, row, clm+1)
	}

	candidates := puzzle[row][clm]
	for d := range candidates {
		if !puzzle.Valid(d, row, clm) {
			continue
		}
		puzzle[row][clm] = map[int]struct{}{d: {}}
		if Sequential(puzzle, row, clm+1) {
			return true
		}
	}
	puzzle[row][clm] = candidates
	return false
}
