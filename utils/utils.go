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

// if V is a pointer (slice, map...) this function doesn't deep copy m
func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	c := map[K]V{}
	for k, v := range m {
		c[k] = v
	}
	return c
}

func IsEmptyIntersection[K comparable, V any](lhs map[K]V, rhs map[K]V) bool {
	m := lhs
	other := rhs
	if len(rhs) < len(m) {
		m = rhs
		other = lhs
	}
	for elem := range m {
		_, ok := other[elem]
		if ok {
			return true
		}
	}
	return false
}
