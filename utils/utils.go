package utils

import (
	"fmt"
	"sync"
	"time"
)

type Solution struct {
	Start            time.Time
	SolverDesign     string
	PuzzleComplexity string
}

func TimeTrack(timer <-chan Solution, wg *sync.WaitGroup) {
	defer wg.Done()
	for tick := range timer {
		fmt.Printf("%s took %f secs to solve %s sudoku\n", tick.SolverDesign, time.Since(tick.Start).Seconds(), tick.PuzzleComplexity)
	}
}
