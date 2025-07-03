package executor

import "github.com/ksysoev/deriv-bot/pkg/core/signal"

type StrategyType int

const (
	StrategyTypeNotSet StrategyType = iota
	StrategyTypeBuy
	StrategyTypeSell
)

type Strategy struct {
	Token        string
	Symbol       string
	Amount       float64
	Type         StrategyType
	Leverage     int
	CheckToOpen  func(tick signal.Tick) bool
	CheckToClose func(tick signal.Tick) bool
}
