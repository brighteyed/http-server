package tracker

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// serverIdleTracker implements IdleTracker and monitors active
// client connections
type serverIdleTracker struct {
	mu     sync.Mutex
	active map[net.Conn]bool
	idle   time.Duration
	timer  *time.Timer
}

// ConnState called when a client connection changes its state
func (t *serverIdleTracker) ConnState(conn net.Conn, state http.ConnState) {
	t.mu.Lock()
	defer t.mu.Unlock()

	oldActive := len(t.active)
	switch state {
	case http.StateNew, http.StateActive, http.StateHijacked:
		t.active[conn] = true
		if oldActive == 0 {
			log.Println("Active state: stopping timer")
			t.timer.Stop()
		}

	case http.StateIdle, http.StateClosed:
		delete(t.active, conn)
		if oldActive > 0 && len(t.active) == 0 {
			log.Println("Inactive state: restarting timer")
			t.timer.Reset(t.idle)
		}
	}
}

// Done signals when idle period expired
func (t *serverIdleTracker) Done() <-chan time.Time {
	return t.timer.C
}
