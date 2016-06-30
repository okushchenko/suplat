package common

import "time"

type Point struct {
	Error     int
	Latency   time.Duration
	Time      time.Time
	Interface string
}
