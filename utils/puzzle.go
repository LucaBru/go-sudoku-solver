package utils

import (
	"fmt"
)

type Digits map[int]struct{}

func (c Digits) Singleton() (int, bool) {
	for key := range c {
		return key, len(c) == 1
	}
	return 0, false
}

func NewPuzzle(board [9][9]int) (puzzle *Puzzle) {
	puzzle = &Puzzle{}
	for i := range 9 {
		for j := range 9 {
			if v := board[i][j]; v == 0 {
				puzzle[i][j] = map[int]struct{}{
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
			} else {
				puzzle[i][j] = map[int]struct{}{v: {}}
			}
		}
	}
	return
}

type Puzzle [9][9]Digits

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

func (p *Puzzle) DeepCopy() *Puzzle {
	dest := *p
	for r := range 9 {
		for c := range 9 {
			dest[r][c] = CopyMap(p[r][c])
		}
	}
	return &dest
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
	return p.clmConstraint(digit, clm) &&
		p.rowConstraint(digit, row) &&
		p.boxConstraint(digit, row, clm)
}

func (p Puzzle) rowConstraint(digit, row int) bool {
	for _, candidates := range p[row] {
		if d, ok := candidates.Singleton(); ok && d == digit {
			return false
		}
	}
	return true
}

func (p Puzzle) clmConstraint(digit, clm int) bool {
	for i := range 9 {
		if d, ok := p[i][clm].Singleton(); ok && d == digit {
			return false
		}
	}
	return true
}

func (p Puzzle) boxConstraint(digit int, row, clm int) bool {
	rowIdx := row / 3 * 3
	clmIdx := clm / 3 * 3
	for i := range 3 {
		for j := range 3 {
			if d, ok := p[rowIdx+i][clmIdx+j].Singleton(); ok && d == digit {
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
