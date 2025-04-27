package internalhttp

import (
	"net/http"
	"time"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/middleware"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(lrw, r)

		logLine := middleware.CommonLogFormat(
			middleware.ExtractIP(r.RemoteAddr),
			r.Method,
			r.URL.Path,
			r.Proto,
			lrw.statusCode,
			time.Since(start),
			r.UserAgent(),
		)

		middleware.WriteLogToFile(logLine)
	})
}
