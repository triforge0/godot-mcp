package events

import "time"

type Event struct {
	Type      string
	Data      any
	Timestamp time.Time
}
