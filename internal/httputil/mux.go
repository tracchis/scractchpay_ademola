package httputil

import (
	"net/http"
	"net/http/pprof"

	"go.opencensus.io/zpages"
)

// NewBaseMux creates a ServeMux with various debugging and status endpoints set up:
// /alive, always 200 OK
// /ready, custom handler provided, `httputil.Ready` can be used for this
// /debug/rpcz, RPC Stats
// /debug/tracez, Trace Spans
// /debug/pprof/, pprof
// /debug/pprof/cmdline,  pprof
// /debug/pprof/profile, pprof
// /debug/pprof/symbol, pprof
// /debug/pprof/trace, pprof
func NewBaseMux(ready http.HandlerFunc) *http.ServeMux {
	mux := http.NewServeMux()

	// /alive always responds with 200 OK
	mux.HandleFunc("/alive", TextHandler(http.StatusOK, "application/json", `"OK"`))

	// /ready is a custom handler, `httputil.Ready` can be used for this
	mux.Handle("/health", ready)

	// zPages exposes various debugging data from OpenCensus
	// endpoints: /debug/rpcz, /debug/tracez
	// more info: https://opencensus.io/zpages/go/
	zpages.Handle(mux, "/debug")

	// pprof allows remote profiling
	// more info: https://golang.org/pkg/net/http/pprof/
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return mux
}
