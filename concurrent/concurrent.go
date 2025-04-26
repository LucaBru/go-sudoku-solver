package concurrent

import (
	"fmt"
	"log"
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

func Solver(puzzle sudoku.Puzzle) sudoku.Puzzle {
	initPrecinctNetwork(puzzle)
	start := make(chan struct{})
	filteredCandidates := make(chan *FilteredCandidates)
	defer close(filteredCandidates)
	detectives := 0
	for i := range 9 {
		for j := range 9 {
			address := Address{row: i, clm: j}
			detective := NewCellDetective(i, j, network[address])
			if ok, value := puzzle[i][j].IsSingleton(); ok {
				detective.notifyColleagues(&Discovery{digit: value, address: address})
				continue
			}
			log.Printf(
				"%s colleagues:\nbox %v\nrow: %v\nclm %v\n",
				address,
				detective.boxDetectivesAddr,
				detective.rowDetectivesAddr,
				detective.clmDetectivesAddr,
			)
			go detective.investigate(start, filteredCandidates)
			detectives++
		}
	}
	close(start)
	for range detectives {
		filtered := <-filteredCandidates
		log.Println(filtered)
		puzzle[filtered.row][filtered.clm] = filtered.candidates
	}
	if puzzle.IsSolved() {
		return puzzle
	}
	fmt.Println("Simplified version of the initial puzzle")
	puzzle.DisplayCandidates()
	return parallel.Solver(puzzle)
}

type FilteredCandidates struct {
	candidates sudoku.Digits
	row        int
	clm        int
}

func (f *FilteredCandidates) String() string {
	return fmt.Sprintf("%s filtered candidates %v\n", Address{row: f.row, clm: f.clm}, f.candidates)
}
