package middleware

import (
	"fmt"
	"net"
	"os"
	"time"
)

func CommonLogFormat(ip, method, path, proto string, status int, latency time.Duration, userAgent string) string {
	date := time.Now().Format("[02/Jan/2006:15:04:05 -0700]")

	if userAgent == "" {
		userAgent = "-"
	} else {
		userAgent = `"` + userAgent + `"`
	}

	return fmt.Sprintf("%s %s %s %s %s %d %d %s\n",
		ip,
		date,
		method,
		path,
		proto,
		status,
		latency.Milliseconds(),
		userAgent,
	)
}

func ExtractIP(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return ip
}

func WriteLogToFile(logLine string) {
	file, err := os.OpenFile("logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %v", err))
	}
	defer file.Close()

	if _, err = file.WriteString(logLine); err != nil {
		panic(fmt.Sprintf("Failed to write to log file: %v", err))
	}
}
