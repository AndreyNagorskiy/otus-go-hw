package internalhttp

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
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

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		date := time.Now().Format("[02/Jan/2006:15:04:05 -0700]")

		userAgent := r.UserAgent()
		if userAgent == "" {
			userAgent = "-"
		} else {
			userAgent = `"` + userAgent + `"`
		}

		latency := time.Since(start).Milliseconds()

		logLine := fmt.Sprintf("%s %s %s %s %s %d %d %s\n",
			ip,
			date,
			r.Method,
			r.URL.Path,
			r.Proto,
			lrw.statusCode,
			latency,
			userAgent,
		)

		writeLogToFile(logLine)
	})
}

func writeLogToFile(logLine string) {
	file, err := os.OpenFile("logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %v", err))
	}
	defer file.Close()

	if _, err = file.WriteString(logLine); err != nil {
		panic(fmt.Sprintf("Failed to write to log file: %v", err))
	}
}
