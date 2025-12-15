//revive:disable
package common

import "time"

type LoggerInterface interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type LogEntry struct {
	IP        string
	Date      time.Time
	Path      string
	Proto     string
	Method    string
	UserAgent string
	Status    int
	Latency   int
}
