package concurrent

import (
	"fmt"
	"maps"
	"sudoku/sudoku"
	"sudoku/utils"
	"time"
)

type Address struct {
	row int
	clm int
}

func (a *Address) rowNeighbors() map[Address]struct{} {
	rowNeighbors := map[Address]struct{}{}
	leftMostClm := a.clm / 3 * 3
	for i := range leftMostClm {
		address := Address{row: a.row, clm: i}
		if network[address] == nil {
			continue
		}
		rowNeighbors[address] = struct{}{}
	}
	for i := range 9 - (leftMostClm + 3) {
		address := Address{row: a.row, clm: leftMostClm + i + 3}
		if network[address] == nil {
			continue
		}
		rowNeighbors[address] = struct{}{}
	}
	return rowNeighbors
}

func (a *Address) clmNeighbors() map[Address]struct{} {
	clmNeighbors := map[Address]struct{}{}
	upperRow := a.row / 3 * 3
	for i := range upperRow {
		address := Address{row: i, clm: a.clm}
		if network[address] == nil {
			continue
		}
		clmNeighbors[address] = struct{}{}
	}
	for i := range 9 - (upperRow + 3) {
		address := Address{row: upperRow + i + 3, clm: a.clm}
		if network[address] == nil {
			continue
		}
		clmNeighbors[address] = struct{}{}
	}
	return clmNeighbors
}

func (a *Address) boxNeighbors() map[Address]struct{} {
	boxNeighbors := map[Address]struct{}{}
	upperRow := a.row / 3 * 3
	leftMostClm := a.clm / 3 * 3
	for i := range 3 {
		for j := range 3 {
			address := Address{row: upperRow + i, clm: leftMostClm + j}
			if network[address] == nil {
				continue
			}
			boxNeighbors[address] = struct{}{}
		}
	}
	delete(boxNeighbors, Address{row: a.row, clm: a.clm})
	return boxNeighbors
}

type CellDetective struct {
	address                Address
	rowDetectivesAddresses map[Address]struct{}
	clmDetectivesAddresses map[Address]struct{}
	boxDetectivesAddresses map[Address]struct{}
	detectivesPerDigit     map[int]map[Address]struct{}
	hints                  <-chan Restriction
}

func NewCellDetective(
	row, clm int,
	hints <-chan Restriction,
	candidates sudoku.Candidates,
) *CellDetective {
	address := Address{row: row, clm: clm}

	rowDetectives := address.rowNeighbors()
	clmDetectives := address.clmNeighbors()
	boxDetectives := address.boxNeighbors()

	detectives := map[Address]struct{}{}
	maps.Copy(detectives, rowDetectives)
	maps.Copy(detectives, clmDetectives)
	maps.Copy(detectives, boxDetectives)
	detectivesPerDigit := map[int]map[Address]struct{}{}
	for i := range 9 {
		detectivesPerDigit[i] = utils.CopyMap(detectives)
	}

	return &CellDetective{
		address:                address,
		rowDetectivesAddresses: rowDetectives,
		clmDetectivesAddresses: clmDetectives,
		boxDetectivesAddresses: boxDetectives,
		detectivesPerDigit:     detectivesPerDigit,
		hints:                  hints,
	}
}

func (d *CellDetective) investigate(start <-chan struct{}, filteredCandidates chan<- *FilteredCandidates) {
	sendFilteredCandidates := func() {
		candidates := map[int]struct{}{}
		for candidate := range d.detectivesPerDigit {
			candidates[candidate] = struct{}{}
		}
		filteredCandidates <- &FilteredCandidates{row: d.address.row, clm: d.address.clm, candidates: candidates}
	}
	defer sendFilteredCandidates()
	<-start
	for {
		discovery := 0
		select {
		case clue := <-d.hints:
			{
				switch clue := clue.(type) {
				case *Discovery:
					discovery = d.onDiscover(clue)
				case *Clue:
					discovery = d.onHint(clue)
				}
			}
		case <-time.After(50 * time.Millisecond):
			return
		}
		if discovery > 0 {
			discovery := &Discovery{digit: discovery,
				row: d.address.row,
				clm: d.address.clm}
			d.notifyColleagues(discovery)
			return
		}
	}
}

func (d *CellDetective) notifyColleagues(msg Restriction) {
	network.notify(d.rowDetectivesAddresses, msg)
	network.notify(d.clmDetectivesAddresses, msg)
	network.notify(d.boxDetectivesAddresses, msg)
}

func (d *CellDetective) giveHint(discovery *Discovery, msg *Clue) {
	network.notify(d.boxDetectivesAddresses, msg)
	if discovery.Row() == d.address.row {
		network.notify(d.clmDetectivesAddresses, msg)
	} else if discovery.Clm() == d.address.clm {
		network.notify(d.rowDetectivesAddresses, msg)
	}
}

func (d *CellDetective) onDiscover(discovery *Discovery) int {
	if d.detectivesPerDigit[discovery.digit] == nil {
		return 0
	}
	fmt.Printf(
		"[%d, %d] recvs discovery digit %d from [%d, %d]\n",
		d.address.row,
		d.address.clm,
		discovery.digit,
		discovery.row,
		discovery.clm,
	)
	delete(d.detectivesPerDigit, discovery.digit)
	if len(d.detectivesPerDigit) == 1 {
		var k int
		for k = range d.detectivesPerDigit {
		}
		fmt.Printf("Detective found spot [%d, %d] = %d\n", d.address.row, d.address.clm, k)
		discovery := &Discovery{digit: discovery.digit, row: d.address.row, clm: d.address.clm}
		d.notifyColleagues(discovery)
		return k
	}
	clue := &Clue{digit: discovery.digit, row: d.address.row, clm: d.address.clm}
	d.giveHint(discovery, clue)
	return 0
}

func (d *CellDetective) onHint(clue *Clue) int {
	if d.detectivesPerDigit[clue.digit] == nil {
		return 0
	}
	fmt.Printf(
		"[%d, %d] recvs clue digit %d from [%d, %d]\n",
		d.address.row,
		d.address.clm,
		clue.digit,
		clue.row,
		clue.clm,
	)
	delete(d.detectivesPerDigit[clue.digit], Address{row: clue.row, clm: clue.clm})
	fmt.Printf(
		"[%d, %d] digit %d suspects %d: %v\nAfter removing suspect [%d, %d]\n",

		d.address.row,
		d.address.clm,
		clue.digit,
		len(d.detectivesPerDigit[clue.digit]),
		d.detectivesPerDigit[clue.digit],
		clue.row,
		clue.clm,
	)
	if d.exclusiveFor(clue.digit) {
		fmt.Printf("[%d, %d] discover %d due to clues\n", d.address.row, d.address.clm, clue.digit)
		return clue.digit
	}
	return 0
}

func (d *CellDetective) exclusiveFor(digit int) bool {
	a := d.detectivesPerDigit[digit]
	return utils.IsEmptyIntersection(a, d.boxDetectivesAddresses) || utils.IsEmptyIntersection(a, d.rowDetectivesAddresses) ||
		utils.IsEmptyIntersection(a, d.clmDetectivesAddresses)
}
