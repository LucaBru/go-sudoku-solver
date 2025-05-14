package concurrent

import "sudoku/utils"

type Network [9][9]chan Rule

type Index struct {
	row uint
	clm uint
}

func NewNetwork(puzzle *utils.Puzzle) (net *Network) {
	net = &Network{}
	for r := range 9 {
		for c := range 9 {
			if _, ok := puzzle[r][c].Singleton(); ok {
				continue
			}
			net[r][c] = make(chan Rule, 10000)
		}
	}
	return
}

func (net *Network) GetMediums(row, clm uint) (rivals *Rivals, rxRule <-chan Rule) {
	rivals = &Rivals{}
	rivals.rows = map[uint]chan<- Rule{}
	rivals.clms = map[uint]chan<- Rule{}
	rivals.boxs = map[Index]chan<- Rule{}
	rxRule = net[row][clm]
	for i := range 9 {
		if net[row][i] != nil {
			rivals.rows[uint(i)] = net[row][i]
		}
		if net[i][clm] != nil {
			rivals.clms[uint(i)] = net[i][clm]
		}
	}

	upperRow := row / 3 * 3
	leftmostClm := clm / 3 * 3
	for rowShift := range 3 {
		row := upperRow + uint(rowShift)
		for clmShift := range 3 {
			clm := leftmostClm + uint(clmShift)
			if net[row][clm] == nil {
				continue
			}
			rivals.boxs[Index{row: row, clm: clm}] = net[row][clm]
		}
	}

	delete(rivals.rows, clm)
	delete(rivals.clms, row)
	delete(rivals.boxs, Index{row: row, clm: clm})
	return
}
