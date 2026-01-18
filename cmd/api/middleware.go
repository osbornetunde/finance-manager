package main

import (
	"net/http"
	"time"
)

// TimeoutMiddleware wraps a handler with a timeout. If the handler doesn't
// complete within the timeout duration, it returns a 503 Service Unavailable.
func (a *API) TimeoutMiddleware(next http.Handler, timeout time.Duration) http.Handler {
	return http.TimeoutHandler(next, timeout, `{"error":"request timeout"}`)
}
