package concurrent

import (
	"log"
	"maps"
	"sudoku/utils"
	"time"
)

type CellDetective struct {
	address           Address
	rowDetectivesAddr map[Address]struct{}
	clmDetectivesAddr map[Address]struct{}
	boxDetectivesAddr map[Address]struct{}
	digitsSuspects    map[int]map[Address]struct{}
	colleaguesHints   <-chan Restriction
}

func NewCellDetective(
	row, clm int,
	hints <-chan Restriction,
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
		detectivesPerDigit[i+1] = utils.CopyMap(detectives)
	}

	return &CellDetective{
		address:           address,
		rowDetectivesAddr: rowDetectives,
		clmDetectivesAddr: clmDetectives,
		boxDetectivesAddr: boxDetectives,
		digitsSuspects:    detectivesPerDigit,
		colleaguesHints:   hints,
	}
}

func (d *CellDetective) investigate(start <-chan struct{}, filteredCandidates chan<- *FilteredCandidates) {
	sendFilteredCandidates := func() {
		candidates := map[int]struct{}{}
		for candidate := range d.digitsSuspects {
			candidates[candidate] = struct{}{}
		}
		filteredCandidates <- &FilteredCandidates{row: d.address.row, clm: d.address.clm, candidates: candidates}
	}
	defer sendFilteredCandidates()
	<-start
	for {
		discovery := 0
		select {
		case clue := <-d.colleaguesHints:
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
			discovery := &Discovery{digit: discovery, address: d.address}
			d.notifyColleagues(discovery)
			return
		}
	}
}

func (d *CellDetective) notifyColleagues(msg Restriction) {
	network.notify(d.rowDetectivesAddr, msg)
	network.notify(d.clmDetectivesAddr, msg)
	network.notify(d.boxDetectivesAddr, msg)
}

func (d *CellDetective) giveHint(discovery *Discovery, msg *Clue) {
	network.notify(d.boxDetectivesAddr, msg)
	if discovery.Row() == d.address.row {
		network.notify(d.clmDetectivesAddr, msg)
		return
	}
	if discovery.Clm() == d.address.clm {
		network.notify(d.rowDetectivesAddr, msg)
	}
}

func (d *CellDetective) onDiscover(discovery *Discovery) int {
	if d.digitsSuspects[discovery.digit] == nil {
		return 0
	}
	hintColleagues := func() {
		clue := &Clue{digit: discovery.digit, address: d.address}
		log.Printf("%s sends %s\n", d.address, clue)
		d.giveHint(discovery, clue)
	}
	defer hintColleagues()
	log.Printf("%s receives %s\n", d.address, discovery)
	delete(d.digitsSuspects, discovery.digit)
	if len(d.digitsSuspects) == 1 {
		var k int
		for k = range d.digitsSuspects {
		}
		log.Printf("%s = %d\n", d.address, k)
		return k
	}
	for digit, suspects := range d.digitsSuspects {
		delete(suspects, discovery.address)
		if len(suspects) == 0 {
			log.Printf("%s = %d\n", d.address, digit)
			return digit
		}
	}
	return 0
}

func (d *CellDetective) onHint(clue *Clue) int {
	if d.digitsSuspects[clue.digit] == nil {
		return 0
	}
	delete(d.digitsSuspects[clue.digit], clue.address)
	log.Printf(
		"%s receives %s. Suspects left: %v\n",
		d.address,
		clue,
		d.digitsSuspects[clue.digit],
	)
	if d.isGuilty(clue.digit) {
		d.digitsSuspects = map[int]map[Address]struct{}{clue.digit: nil}
		log.Printf("%s = %d thanks to clues\n", d.address, clue.digit)
		return clue.digit
	}
	return 0
}

func (d *CellDetective) isGuilty(digit int) bool {
	return utils.IsEmptyIntersection(d.digitsSuspects[digit], d.boxDetectivesAddr) ||
		utils.IsEmptyIntersection(d.digitsSuspects[digit], d.rowDetectivesAddr) ||
		utils.IsEmptyIntersection(d.digitsSuspects[digit], d.clmDetectivesAddr)
}
