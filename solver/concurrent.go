package solver

import (
	"log"
	"sudoku/concurrent"
	"sudoku/utils"
	"sync"
)

func Concurrent(puzzle *utils.Puzzle) bool {
	log.Println("start")
	net := concurrent.NewNetwork(puzzle)
	start := make(chan struct{})
	rxDigit := make(chan concurrent.Rule)
	wg := &sync.WaitGroup{}
	for i := range 9 {
		for j := range 9 {
			rivals, rxRule := net.GetMediums(uint(i), uint(j))
			if digit, ok := puzzle[i][j].Singleton(); ok {
				rivals.NotifyAll(concurrent.Rule{
					Digit: uint(digit),
					Row:   uint(i),
					Clm:   uint(j),
					Kind:  concurrent.Discovery,
				})
				continue
			}
			cell := concurrent.NewCell(uint(i), uint(j), rivals, rxRule)
			wg.Add(1)
			go cell.Investigate(start, rxDigit, wg)
		}
	}

	go func() {
		wg.Wait()
		close(rxDigit)
	}()

	close(start)

	for digit := range rxDigit {
		log.Printf("Spot digit %d at %d %d\n", digit.Digit, digit.Row, digit.Clm)
		puzzle[digit.Row][digit.Clm] = map[int]struct{}{int(digit.Digit): {}}
	}
	return Parallel(puzzle)
}
