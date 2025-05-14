package concurrent

import (
	"log"
	"sudoku/utils"
)

type Digit = uint

type Rivals struct {
	rows map[uint]chan<- Rule
	clms map[uint]chan<- Rule
	boxs map[Index]chan<- Rule
}

func notify(ch chan<- Rule, msg Rule) {
	if ch != nil {
		ch <- msg
	}
}

func (r *Rivals) NotifyAll(msg Rule) {
	log.Printf("%d %d notify %v %v %v\n", msg.Row, msg.Clm, r.rows, r.clms, r.boxs)
	for _, ch := range r.rows {
		notify(ch, msg)
	}
	for _, ch := range r.clms {
		notify(ch, msg)
	}
	for _, ch := range r.boxs {
		notify(ch, msg)
	}
}

func (r *Rivals) clone() *Rivals {
	return &Rivals{
		rows: utils.CopyMap(r.rows),
		clms: utils.CopyMap(r.clms),
		boxs: utils.CopyMap(r.boxs),
	}
}

func (r *Rivals) deleteDueTo(rule Rule, row uint, clm uint) {
	delete(r.boxs, Index{row: rule.Row, clm: rule.Clm})
	if rule.Row == row {
		delete(r.rows, rule.Clm)
	}
	if rule.Clm == clm {
		delete(r.clms, rule.Row)
	}
}

func (r *Rivals) Len() (int, int, int) {
	return len(r.rows), len(r.clms), len(r.boxs)
}

func (r *Rivals) Lose() bool {
	return len(r.rows) == 0 || len(r.clms) == 0 || len(r.boxs) == 0
}
