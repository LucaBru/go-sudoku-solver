package concurrent

type RuleType int

const (
	Discovery RuleType = iota
	Clue
)

type Rule struct {
	Digit uint
	Row   uint
	Clm   uint
	Kind  RuleType
}
