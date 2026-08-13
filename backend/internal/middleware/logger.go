package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf("[%s] %s | Status: %d | Duration: %v", r.Method, r.URL.Path, rw.statusCode, duration)

		if rw.statusCode >= 500 {
			log.Printf("[SERVER ERROR] %s %s returned status %d", r.Method, r.URL.Path, rw.statusCode)
		}
	})
}
