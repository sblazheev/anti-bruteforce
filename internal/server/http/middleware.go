package internalhttp

import (
	"net/http"
	"strings"
	"time"

	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common" //nolint:depguard
)

func loggingMiddleware(next http.Handler, logger common.LoggerInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := time.Now()
		rw := NewStatusResponseWriter(w)
		next.ServeHTTP(rw, r)

		l := time.Since(s)
		logger.Info("REQUEST API",
			"data",
			struct {
				IP        string
				Date      time.Time
				Path      string
				Proto     string
				Method    string
				UserAgent string
				Status    int
				Latency   int
			}{
				IP:        strings.Split(r.RemoteAddr, ":")[0],
				Date:      s,
				Path:      r.URL.Path,
				Proto:     r.Proto,
				Method:    r.Method,
				UserAgent: r.UserAgent(),
				Status:    rw.statusCode,
				Latency:   int(l.Milliseconds()),
			})
	})
}

func errorJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextW := NewStatusResponseWriter(w)
		next.ServeHTTP(nextW, r)
	})
}
