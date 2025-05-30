package main

import (
	"log"
	"os"

	"sudoku/solver"
	"sudoku/utils"
	"sync"
	"time"
)

func init() {
	os.Remove("solver.log")
	logFile, _ := os.OpenFile("solver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	boards := map[string][9][9]int{
		/* "naiveBoard": [9][9]int{
			{3, 1, 2, 6, 0, 5, 4, 0, 0},
			{6, 0, 4, 2, 1, 0, 0, 8, 3},
			{9, 0, 8, 0, 3, 0, 0, 2, 0},
			{2, 4, 7, 5, 6, 0, 0, 3, 0},
			{8, 6, 0, 1, 0, 0, 0, 0, 0},
			{0, 0, 5, 3, 0, 2, 6, 7, 0},
			{0, 8, 0, 0, 0, 0, 0, 0, 4},
			{0, 3, 0, 0, 0, 0, 7, 6, 2},
			{5, 0, 0, 0, 7, 0, 8, 0, 9},
		}, */

		"testSimpTechniques": [9][9]int{
			{0, 0, 0, 1, 0, 4, 0, 0, 0},
			{0, 0, 1, 0, 0, 0, 9, 0, 0},
			{0, 9, 0, 7, 0, 3, 0, 6, 0},
			{8, 0, 7, 0, 0, 0, 1, 0, 6},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{3, 0, 4, 0, 0, 0, 5, 0, 9},
			{0, 5, 0, 4, 0, 2, 0, 3, 0},
			{0, 0, 8, 0, 0, 0, 6, 0, 0},
			{0, 0, 0, 8, 0, 6, 0, 0, 0},
		},

		/* "quiteHardBoard": [9][9]int{
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 3, 0, 8, 5},
			{0, 0, 1, 0, 2, 0, 0, 0, 0},
			{0, 0, 0, 5, 0, 7, 0, 0, 0},
			{0, 0, 4, 0, 0, 0, 1, 0, 0},
			{0, 9, 0, 0, 0, 0, 0, 0, 0},
			{5, 0, 0, 0, 0, 0, 0, 7, 3},
			{0, 0, 2, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 4, 0, 0, 0, 9},
		}, */

		/* "hardestBoard": [9][9]int{
			{9, 0, 0, 8, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 5, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 2, 0, 0, 1, 0, 0, 0, 3},
			{0, 1, 0, 0, 0, 0, 0, 6, 0},
			{0, 0, 0, 4, 0, 0, 0, 7, 0},
			{7, 0, 8, 6, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 3, 0, 1, 0, 0},
			{4, 0, 0, 0, 0, 0, 2, 0, 0},
		}, */
	}

	timer := make(chan utils.Solution)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go utils.TimeTrack(timer, &wg)

	for key, board := range boards {
		sPuzzle := utils.NewPuzzle(board)
		msg := utils.Solution{Start: time.Now(), SolverDesign: "sequential", PuzzleComplexity: key}
		solver.Sequential(sPuzzle, 0, 0)
		msg.Sol = sPuzzle
		timer <- msg

		pPuzzle := utils.NewPuzzle(board)
		msg.SolverDesign = "parallel"
		msg.Start = time.Now()
		solver.Parallel(pPuzzle)
		msg.Sol = pPuzzle
		timer <- msg

		cPuzzle := utils.NewPuzzle(board)
		msg.SolverDesign = "concurrent"
		msg.Start = time.Now()
		solver.Concurrent(cPuzzle)
		msg.Sol = cPuzzle
		timer <- msg
	}
	close(timer)
	wg.Wait()
}
