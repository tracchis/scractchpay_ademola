package httputil

import (
	"net/http"
	"sync"
)

// Ready is a structure that provides safe readiness handling.
//
// A Ready is in one of two states, it is either "unready" (its initial state),
// or "ready" (serving traffic).
//
// Using Handler(), a Ready returns a wrapped http.Handler that will call either
// the internally wrapped unready http.Handler specified in NewReady when "unready",
// or call the given http.Handler when "ready".
type Ready struct {
	unreadyHandler http.Handler

	mu sync.RWMutex
	ch chan struct{}
}

// NewReady returns a middleware that provides reliable and safe readiness handling.
func NewReady(unready http.Handler) *Ready {
	return &Ready{
		unreadyHandler: unready,

		ch: make(chan struct{}),
	}
}

// State returns the current state of the readiness handler.
func (r *Ready) State() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	select {
	case <-r.ch:
		return "unready"
	default:
		return "ready"
	}
}

// Notify returns a receive-only channel that will be closed when the Ready enters the "ready" state.
func (r *Ready) Notify() <-chan struct{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.ch
}

// Ready transitions the Ready handler into the "ready" state.
//
// Ready then unlocks all Handlers from this Ready and will begin performing their regular functions.
func (r *Ready) Ready() {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-r.ch:
		return
	default:
	}

	close(r.ch)
}

// Unready transitions the Ready handler into the "unready" state.
//
// Unready then locks all Handlers from this Ready and will begin returning the built in unready handler.
func (r *Ready) Unready() {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-r.ch:
	default:
		return
	}

	r.ch = make(chan struct{})
}

// Handler returns an http.HandlerFunc that will call the internal unready handler from Ready when it is "unready".
// After being marked "ready", the this handler will call the given ready handler.
func (r *Ready) Handler(ready http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		select {
		case <-r.Notify():
			ready.ServeHTTP(rw, req)
		default:
			r.unreadyHandler.ServeHTTP(rw, req)
		}
	}
}
