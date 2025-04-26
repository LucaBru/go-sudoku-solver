package concurrent

type Restriction interface {
	Digit() int
	Row() int
	Clm() int
}

type Discovery struct {
	digit int
	row   int
	clm   int
}

func (r *Discovery) Digit() int {
	return r.digit
}

func (r *Discovery) Row() int {
	return r.row
}

func (r *Discovery) Clm() int {
	return r.clm
}

type Clue struct {
	digit int
	row   int
	clm   int
}

func (r *Clue) Digit() int {
	return r.digit
}

func (r *Clue) Row() int {
	return r.row
}

func (r *Clue) Clm() int {
	return r.clm
}
