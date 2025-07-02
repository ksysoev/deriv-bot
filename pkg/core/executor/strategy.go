package executor

type StrategyType int

const (
	StrategyTypeNotSet StrategyType = iota
	StrategyTypeBuy
	StrategyTypeSell
)

type Strategy struct {
	Token        string
	Symbol       string
	Amount       string
	Type         StrategyType
	Leverage     int
	CheckToOpen  func(any) bool
	CheckToClose func(any) bool
}
