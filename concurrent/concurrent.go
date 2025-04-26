package concurrent

import (
	"fmt"
	"sudoku/parallel"
	"sudoku/sudoku"
)

type PrecinctNetwork map[Address]chan Restriction

var network PrecinctNetwork

func (n PrecinctNetwork) notify(addresses map[Address]struct{}, msg Restriction) {
	for address := range addresses {
		n[address] <- msg
	}
}

func initPrecinctNetwork(puzzle sudoku.Puzzle) {
	m := map[Address]chan Restriction{}
	for i, row := range puzzle {
		for j, candidates := range row {
			var ch chan Restriction
			if ok, _ := candidates.IsSingleton(); !ok {
				ch = make(chan Restriction, 5000)
			}
			m[Address{row: i, clm: j}] = ch
		}
	}
	network = m
}

func Solve(puzzle sudoku.Puzzle) sudoku.Puzzle {
	initPrecinctNetwork(puzzle)
	start := make(chan struct{})
	filteredCandidates := make(chan *FilteredCandidates)
	detectives := 0
	for i := range 9 {
		for j := range 9 {
			detective := NewCellDetective(i, j, network[Address{row: i, clm: j}], puzzle[i][j])
			fmt.Printf(
				"[%d, %d] colleagues:\nbox %v\nrow: %v\nclm %v\n",
				i,
				j,
				detective.boxDetectivesAddresses,
				detective.rowDetectivesAddresses,
				detective.clmDetectivesAddresses,
			)
			if ok, value := puzzle[i][j].IsSingleton(); ok {
				detective.notifyColleagues(&Discovery{digit: value, row: i, clm: j})
				continue
			}
			go detective.investigate(start, filteredCandidates)
			detectives++
		}
	}
	close(start)
	for range detectives {
		filtered := <-filteredCandidates
		puzzle[filtered.row][filtered.clm] = filtered.candidates
	}
	close(filteredCandidates)
	if puzzle.IsSolved() {
		return puzzle
	}
	return parallel.Solver(puzzle)
}

type FilteredCandidates struct {
	candidates sudoku.Candidates
	row        int
	clm        int
}
