package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)
		slog.Info(fmt.Sprintf("[%s] %s %s %d %s", strings.Split(r.RemoteAddr, ":")[0], r.Method, r.URL, recorder.statusCode, r.UserAgent()))
	})
}
