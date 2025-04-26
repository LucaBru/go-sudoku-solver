package concurrent

import "fmt"

type Restriction interface {
	Digit() int
	Row() int
	Clm() int
	String() string
}

type Discovery struct {
	digit   int
	address Address
}

func (r *Discovery) Digit() int {
	return r.digit
}

func (r *Discovery) Row() int {
	return r.address.row
}

func (r *Discovery) Clm() int {
	return r.address.clm
}

func (r *Discovery) String() string {
	return fmt.Sprintf("discovery digit %d from %s", r.digit, r.address)
}

type Clue struct {
	digit   int
	address Address
}

func (r *Clue) Digit() int {
	return r.digit
}

func (r *Clue) Row() int {
	return r.address.row
}

func (r *Clue) Clm() int {
	return r.address.clm
}

func (r *Clue) String() string {
	return fmt.Sprintf("clue digit %d from %s", r.digit, r.address)
}
