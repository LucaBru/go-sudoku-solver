package concurrent

import (
	"log"
	"sync"
	"time"
)

type Cell struct {
	row    uint
	clm    uint
	digits map[Digit]*Rivals
	rxRule <-chan Rule
}

func NewCell(
	row, clm uint,
	rivals *Rivals,
	hints <-chan Rule,
) *Cell {
	digits := map[Digit]*Rivals{}
	for i := range 9 {
		digits[uint(i+1)] = rivals.clone()
	}
	return &Cell{
		row:    row,
		clm:    clm,
		digits: digits,
		rxRule: hints,
	}
}

func (d *Cell) Investigate(start <-chan struct{}, txSol chan<- Rule, wg *sync.WaitGroup) {
	defer wg.Done()
	<-start
	for {
		select {
		case rule := <-d.rxRule:
			log.Printf("%d %d has received %+v\n", d.row, d.clm, rule)
			{
				var digit uint
				var ok bool
				switch rule.Kind {
				case Discovery:
					digit, ok = d.onDiscovery(rule)
				case Clue:
					digit, ok = d.onHint(rule)
				}
				if ok {
					discovery := Rule{
						Digit: digit,
						Row:   d.row,
						Clm:   d.clm,
						Kind:  Discovery,
					}
					d.digits[digit].NotifyAll(discovery)
					txSol <- discovery
					return
				}
			}
		case <-time.After(50 * time.Millisecond):
			log.Println("Timeout expired")
			return
		}
	}
}

func (d *Cell) onDiscovery(discovery Rule) (uint, bool) {
	rivals, ok := d.digits[discovery.Digit]
	if !ok {
		return 0, false
	}
	delete(d.digits, discovery.Digit)
	if len(d.digits) == 1 {
		var digit Digit
		for k := range d.digits {
			digit = k
			break
		}
		return digit, true
	}

	for digit, rivals := range d.digits {
		rivals.deleteDueTo(discovery, d.row, d.clm)
		if digit == discovery.Digit {

			log.Printf("%d %d updates digit %d rivals due to %+v\n%v\n%v\n%v\n", d.row, d.clm, digit, discovery, rivals.rows, rivals.clms, rivals.boxs)
		}
		if rivals.Lose() {
			return digit, true
		}
	}

	rivals.NotifyAll(Rule{
		Digit: discovery.Digit,
		Row:   d.row,
		Clm:   d.clm,
		Kind:  Clue,
	})
	return 0, false
}

func (d *Cell) onHint(clue Rule) (uint, bool) {
	digit := clue.Digit
	rivals := d.digits[digit]
	if rivals == nil {
		return 0, false
	}

	rivals.deleteDueTo(clue, d.row, d.clm)
	log.Printf("%d %d updates digit %d rivals due to %+v\n%v\n%v\n%v\n", d.row, d.clm, digit, clue, rivals.rows, rivals.clms, rivals.boxs)

	if rivals.Lose() {
		return digit, true
	}

	return 0, false
}
