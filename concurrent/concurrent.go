package concurrent

import (
	"fmt"
	"sudoku/parallel"
	"sudoku/sudoku"
	"sync"
	"time"
)

/*
each go routine must have:
	- 22 sending channels
	- 22 receiving channels
*/

// any should be a chan (recv or sending of any type actually)
type Neighborhood struct {
	row []chan<- *Restriction
	clm []chan<- *Restriction
	box []chan<- *Restriction
}

func (n *Neighborhood) notify(r *Restriction) {
	for _, ch := range n.row {
		ch <- r
	}

	for _, ch := range n.clm {
		ch <- r
	}

	for _, ch := range n.box {
		ch <- r
	}
}

type Restriction struct {
	digit int
	row   int
	clm   int
}

type News struct {
	Restriction
}

type Hint struct {
	Restriction
}

type ChannelsHandler [9][9]chan *Restriction

func newChannelsHandler(puzzle sudoku.Puzzle) ChannelsHandler {
	matrix := [9][9]chan *Restriction{}
	for i := range 9 {
		for j := range 9 {
			if len(puzzle[i][j]) == 1 {
				continue
			}
			matrix[i][j] = make(chan *Restriction, 22)
		}
	}
	return matrix
}

func (h ChannelsHandler) getChannels(row, clm int) (*Neighborhood, <-chan *Restriction, error) {
	if row > 9 || row < 0 || clm > 9 || clm < 0 {
		return nil, nil, fmt.Errorf("invalid row or clm parameters")
	}

	rowChs := []chan<- *Restriction{}
	for i := range 9 {
		if i == clm || h[row][i] == nil {
			continue
		}
		rowChs = append(rowChs, h[row][i])
	}

	clmChs := []chan<- *Restriction{}
	for i := range 9 {
		if i == row || h[i][clm] == nil {
			continue
		}
		clmChs = append(clmChs, h[i][clm])
	}

	boxChs := []chan<- *Restriction{}
	upperBoxRow := row / 3 * 3
	leftmostBoxClm := clm / 3 * 3
	for i := range 3 {
		for j := range 3 {
			k := upperBoxRow + i
			l := leftmostBoxClm + j
			if row != k && clm == l || h[k][l] == nil {
				continue
			}
			boxChs = append(boxChs, h[k][l])
		}
	}

	return &Neighborhood{row: rowChs, clm: clmChs, box: boxChs}, h[row][clm], nil
}

func Solver(puzzle sudoku.Puzzle) sudoku.Puzzle {
	start := time.Now()
	done := make(chan struct{})
	channels := newChannelsHandler(puzzle)
	wg := &sync.WaitGroup{}
	candidatesRemovedByClues := 0
	for i, row := range puzzle {
		for j, cell := range row {
			neighborhood, recvRestriction, err := channels.getChannels(i, j)
			if ok, d := cell.IsSingleton(); ok {
				candidatesRemovedByClues += len(neighborhood.row) + len(neighborhood.clm) + len(neighborhood.box)
				neighborhood.notify(&Restriction{digit: d, row: i, clm: j})
				continue
			}
			wg.Add(1)
			if err != nil {
				close(done)
				return nil
			}
			go cellHandler(puzzle[i][j], i, j, done, recvRestriction, neighborhood, wg)
		}
	}
	wg.Wait()
	fmt.Printf("Concurrent simplification took %f secs\n", time.Since(start).Seconds())
	close(done)
	CandidatesRemoved.decrement(candidatesRemovedByClues)
	if puzzle.IsSolved() {
		return puzzle
	}
	return parallel.Solver(puzzle)
}

func cellHandler(candidates sudoku.Candidates, row, clm int, done <-chan struct{}, recvRestriction <-chan *Restriction, neighborhood *Neighborhood, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		case r := <-recvRestriction:
			{
				CandidatesRemoved.increment()
				if r.digit < 1 || r.digit > 9 {
					break
				}
				delete(candidates, r.digit)
				if ok, d := candidates.IsSingleton(); ok {
					neighborhood.notify(&Restriction{digit: d, row: row, clm: clm})
					return
				}
			}
		case <-time.After(100 * time.Millisecond):
			{
				return
			}
		}
	}
}

type Counter struct {
	sync.Mutex
	count int
}

func (c *Counter) increment() {
	c.Lock()
	defer c.Unlock()
	c.count++
}

func (c *Counter) value() int {
	c.Lock()
	defer c.Unlock()
	return c.count
}

func (c *Counter) decrement(v int) {
	c.Lock()
	defer c.Unlock()
	c.count -= v
}

var CandidatesRemoved Counter = Counter{count: 0}

/*
This technique doesn't give you the result for all sudoku, indeed for some one you just have to try and pray
But in low time I will get a significant simplification of the board, so if no solution was found, I can apply the parallel classic sudoku solver and solve it.
TODO:

nevertheless it can be improved removing more digits from each cell.

Whenever a cell recv a message:
- from row neighbor -> notify others row in the neighborhood
- from clm -> notify others clm in the neighborhood

Need to distinguish between different kind of messages
- news as always
- hints: remove the neighbor from the list of neighbors that can have the spec digit. If I am the only one remaining, then I conclude I will have that number.
So I can notify the others and conclude
*/
