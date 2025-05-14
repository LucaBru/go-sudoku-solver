package solver

import (
	"sudoku/utils"
	"sync"
)

func Parallel(puzzle *utils.Puzzle) bool {
	rxSol := make(chan *utils.Puzzle)
	row, clm := densestCell(puzzle)
	wg := sync.WaitGroup{}
	for d := range puzzle[row][clm] {
		simplerPuzzle := puzzle.DeepCopy()
		simplerPuzzle[row][clm] = map[int]struct{}{d: {}}
		wg.Add(1)
		go func() {
			if !Sequential(simplerPuzzle, 0, 0) {
				return
			}
			rxSol <- simplerPuzzle
		}()
	}

	go func() {
		wg.Wait()
		close(rxSol)
	}()

	sol, ok := <-rxSol
	*puzzle = *sol
	return ok
}

func densestCell(puzzle *utils.Puzzle) (row, clm uint) {
	density := 0
	for row = range 9 {
		for clm = range 9 {
			if d := len(puzzle[row][clm]); d > density {
				density = d
				if density == 9 {
					return
				}
			}
		}
	}
	return
}
