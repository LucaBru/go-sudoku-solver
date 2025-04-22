package sudoku

import "fmt"

type Candidates map[int]struct{}

func (c Candidates) IsSingleton() (bool, int) {
	for key := range c {
		return len(c) == 1, key
	}
	return false, 0
}

func NewPuzzle(board [][]int) [][]Candidates {
	puzzle := make([][]Candidates, 9)
	for i := range 9 {
		row := make([]Candidates, 9)
		for j := range 9 {
			if v := board[i][j]; v == 0 {
				row[j] = map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}, 5: {}, 6: {}, 7: {}, 8: {}, 9: {}}
			} else {
				row[j] = map[int]struct{}{v: {}}
			}
		}
		puzzle[i] = row
	}
	return puzzle
}

type Puzzle [][]Candidates

func (p Puzzle) Display() {
	for i, row := range p {
		fmt.Printf("Row %d:\n", i)
		for j, candidate := range row {
			fmt.Printf("  Column %d: %v\n", j, candidate)
		}
	}
}

func (p Puzzle) Valid(digit int, row, clm int) bool {
	return p.satisfyClmConstraint(digit, clm) && p.satisfyRowConstraint(digit, row) && p.satisfyBoxConstraint(digit, row, clm)
}

func (p Puzzle) DeepCopy() Puzzle {
	dest := make([][]Candidates, 9)
	for i := range 9 {
		row := make([]Candidates, 9)
		copy(row, p[i])
		dest[i] = row
	}
	return dest
}

func (p Puzzle) satisfyRowConstraint(digit, row int) bool {
	for _, candidates := range p[row] {
		if ok, d := candidates.IsSingleton(); ok && d == digit {
			return false
		}
	}
	return true
}

func (p Puzzle) satisfyClmConstraint(digit, clm int) bool {
	for i := range 9 {
		if ok, d := p[i][clm].IsSingleton(); ok && d == digit {
			return false
		}
	}
	return true
}

func (p Puzzle) satisfyBoxConstraint(digit int, row, clm int) bool {
	rowIdx := row / 3 * 3
	clmIdx := clm / 3 * 3
	for i := range 3 {
		for j := range 3 {
			if ok, d := p[rowIdx+i][clmIdx+j].IsSingleton(); ok && d == digit {
				return false
			}
		}
	}
	return true
}

func (p Puzzle) IsSolved() bool {
	for i := range 9 {
		for j := range 9 {
			if len(p[i][j]) != 1 {
				return false
			}
		}
	}
	return true
}
