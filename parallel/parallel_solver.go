package parallel

import (
	"sudoku/sequential"
	"sudoku/sudoku"
)

/*
when a go routine terminate -> a number cannot be in a cell
all other go routine should be notified, and:
  - if it have already tried that number in that cell do nothing,
  - if they are trying that number stop the dfs and increment the number
  - if it has not tried the number yet, exclude its

the main goroutine should be notified too, and it won't start a go routine with that number
*/

/*
in the solution each cell has a number between 1 and 9, so the number of goroutine depends on how I create sub problems
Trying to guess just one cell ==> 9 goroutine: in a 8 core machine should be ok

Sudoku puzzle has one unique solution (17 clues required)
*/

func SolveSubPuzzle(puzzle sudoku.Puzzle, sendSol chan<- sudoku.Puzzle) {
	sol := sequential.Solver(puzzle, sudoku.Pos{Row: 0, Clm: 0})
	if sol == nil {
		return
	}
	sendSol <- sol
}

func Solver(puzzle sudoku.Puzzle) sudoku.Puzzle {
	recvSol := make(chan sudoku.Puzzle)
	pos := &sudoku.Pos{Row: 0, Clm: 0}
	availableDigits := []int{}
	for i := range 9 {
		for j := range 9 {
			if ok, s := puzzle[i][j].IsSingleton(); !ok && len(s) > len(availableDigits) {
				availableDigits = s
				pos.Row = i
				pos.Clm = j
			}
		}
	}

	for _, d := range availableDigits {
		subPuzzle := puzzle.DeepCopy()
		subPuzzle[pos.Row][pos.Clm].SetSingleton(d)
		go SolveSubPuzzle(subPuzzle, recvSol)
	}
	sol := <-recvSol
	return sol
}
