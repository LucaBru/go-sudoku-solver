package sudoku

import (
	"fmt"
)

type Digits map[int]struct{}

func (c Digits) IsSingleton() (bool, int) {
	for key := range c {
		return len(c) == 1, key
	}
	return false, 0
}

func NewPuzzle(board [][]int) [][]Digits {
	puzzle := make([][]Digits, 9)
	for i := range 9 {
		row := make([]Digits, 9)
		for j := range 9 {
			if v := board[i][j]; v == 0 {
				myMpa := map[int]struct{}{
					1: {},
					2: {},
					3: {},
					4: {},
					5: {},
					6: {},
					7: {},
					8: {},
					9: {},
				}
				row[j] = myMpa
			} else {
				row[j] = map[int]struct{}{v: {}}
			}
		}
		puzzle[i] = row
	}
	return puzzle
}

type Puzzle [][]Digits

func (p Puzzle) Display() {
	fmt.Println("Sudoku")
	for _, row := range p {
		for _, candidate := range row {
			for k := range candidate {
				fmt.Printf(" %d |", k)
				break
			}
		}
		fmt.Printf("\n")
	}
}

func (p Puzzle) DisplayCandidates() {
	fmt.Println("Sudoku")
	for _, row := range p {
		for _, candidate := range row {
			row := "["
			for k := range candidate {
				row += fmt.Sprintf(" %d", k)
			}
			row += " ] | "
			fmt.Printf("%s", row)
		}
		fmt.Printf("\n")
	}
}

func (p Puzzle) Valid(digit int, row, clm int) bool {
	return p.satisfyClmConstraint(digit, clm) && p.satisfyRowConstraint(digit, row) &&
		p.satisfyBoxConstraint(digit, row, clm)
}

func (p Puzzle) DeepCopy() Puzzle {
	dest := make([][]Digits, 9)
	for i := range 9 {
		row := make([]Digits, 9)
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
				fmt.Printf("%d %d length %d\n", i, j, len(p[i][j]))
				return false
			}
		}
	}
	return true
}
