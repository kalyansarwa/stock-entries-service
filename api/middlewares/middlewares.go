package middlewares

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {

			requestID := w.Header().Get("X-Request-Id")
			if requestID == "" {
				requestID = "unknown"
			}
			log.Printf("[%s] %s %s (From: %s, Agent: %s) Time to Respond: %s", requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), time.Since(start))
		}(time.Now())
		next.ServeHTTP(w, r)
	})
}

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = strconv.FormatInt(time.Now().UnixNano(), 36)
		}
		w.Header().Set("X-Request-Id", requestID)
		next.ServeHTTP(w, r)
	})
}
