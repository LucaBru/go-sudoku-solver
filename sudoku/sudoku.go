package sudoku

import (
	"fmt"
	"strings"
)


type Digits [9]bool

// if I use []int rather than [9]bool this function has O(1) cost (by using array length). Don't actually really needed
func (d *Digits) IsSingleton() (bool, []int) {
	ok := []int{}
	for i, d := range d {
		if d {
			ok = append(ok, i+1)
		}
	}
	return len(ok) == 1, ok
}

// less performance O(9)
func (d *Digits) Exclude(v int) {
	if v < 1 || v > 9 {
		return
	}
	d[v-1] = false
}

// more performance O(1)
func Singleton(d int) Digits {
	digits := [9]bool{}
	digits[d-1] = true
	return digits
}

// more performance O(1)
func (d *Digits) SetSingleton(v int) {
	for i := range 9 {
		d[i] = false
	}
	d[v-1] = true
}

func NewDigits() Digits {
	digits := [9]bool{}
	for i := range digits {
		digits[i] = true
	}
	return digits
}

type Puzzle [][]Digits

func (p Puzzle) DeepCopy() Puzzle {
	dest := make([][]Digits, 9)
	for i := range 9 {
		row := make([]Digits, 9)
		copy(row, p[i])
		dest[i] = row
	}
	return dest
}

func (p Puzzle) IsSolved() bool {
	for i := range 9 {
		for j := range 9 {
			if ok, _ := p[i][j].IsSingleton(); !ok {
				return false
			}
		}
	}
	return true
}

func NewPuzzle(board [][]int) [][]Digits {
	puzzle := make([][]Digits, 9)
	for i := range 9 {
		row := make([]Digits, 9)
		for j := range 9 {
			if v := board[i][j]; v != 0 {
				row[j] = Singleton(v)
			} else {
				row[j] = NewDigits()
			}
		}
		puzzle[i] = row
	}
	return puzzle
}

func (p Puzzle) Valid(digit int, pos Pos) bool {
	return p.satisfyClmConstraint(digit, pos.Clm) && p.satisfyRowConstraint(digit, pos.Row) && p.satisfyBoxConstraint(digit, pos)
}

func (p Puzzle) PrettyPrint() {
	for _, row := range p {
		var line []string
		for _, cell := range row {
			var digits []string
			for i, val := range cell {
				if val {
					digits = append(digits, fmt.Sprintf("%d", i+1))
				}
			}
			if len(digits) == 0 {
				line = append(line, ".")
			} else {
				line = append(line, strings.Join(digits, ""))
			}
		}
		fmt.Println(strings.Join(line, " | "))
	}
}

type Pos struct {
	Row int
	Clm int
}

func (p *Pos) Next() *Pos {
	clm := p.Clm + 1
	row := p.Row
	if clm == 9 {
		clm = 0
		row += 1
	}

	return &Pos{
		Row: row,
		Clm: clm,
	}
}

func (p Puzzle) satisfyRowConstraint(digit, row int) bool {
	for _, cell := range p[row] {
		if ok, s := cell.IsSingleton(); ok && s[0] == digit {
			return false
		}
	}
	return true
}

func (p Puzzle) satisfyClmConstraint(digit, clm int) bool {
	for i := range 9 {
		if ok, s := p[i][clm].IsSingleton(); ok && s[0] == digit {
			return false
		}
	}
	return true
}

func (p Puzzle) satisfyBoxConstraint(digit int, pos Pos) bool {
	rowIdx := pos.Row / 3 * 3
	clmIdx := pos.Clm / 3 * 3
	for i := range 3 {
		for j := range 3 {
			if ok, s := p[rowIdx+i][clmIdx+j].IsSingleton(); ok && s[0] == digit {
				return false
			}
		}
	}
	return true
}
