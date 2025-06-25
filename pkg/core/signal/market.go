package signal

import (
	"time"
)

type Tick struct {
	Quote float64
	Ask   float64
	Bid   float64
	Time  time.Time
}
