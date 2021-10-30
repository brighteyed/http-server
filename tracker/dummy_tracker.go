package tracker

import (
	"net"
	"net/http"
	"time"
)

// dummyIdleTracker implements IdleTracker and
// does nothing
type dummyIdleTracker struct {
}

// ConnState called when a client connection changes its state
func (*dummyIdleTracker) ConnState(conn net.Conn, state http.ConnState) {
}

// Done signals when idle period expired. Dummy tracker never signals
func (t *dummyIdleTracker) Done() <-chan time.Time {
	return make(<-chan time.Time)
}
