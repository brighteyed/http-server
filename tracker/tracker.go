package tracker

import (
	"net"
	"net/http"
	"time"
)

// IdleTracker is implemented by idle state trackers
type IdleTracker interface {
	ConnState(conn net.Conn, state http.ConnState)
	Done() <-chan time.Time
}

// NewIdleTracker creates idle tracker. If idleDuration equals to zero,
// then dummy tracker created
func NewIdleTracker(idleDuration uint) IdleTracker {
	if idleDuration != 0 {
		return &serverIdleTracker{
			timer:  time.NewTimer(time.Duration(idleDuration) * time.Second),
			idle:   time.Duration(idleDuration) * time.Second,
			active: make(map[net.Conn]bool),
		}
	}

	return &dummyIdleTracker{}
}
