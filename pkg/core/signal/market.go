package signal

import (
	"time"
)

type Tick struct {
	Time  time.Time
	Quote float64
	Ask   float64
	Bid   float64
}
