package logger

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-universal/console"
)

// logData represents a key-value pair of metadata for a log entry.
type logData struct {
	key   string
	value any
}

// log represents a log entry with a level, timestamp, and associated data.
type log struct {
	level     Level
	timestamp time.Time
	data      []logData
}

// LogOptions defines a function type for modifying a Log instance.
type LogOptions func(*log)

// With adds extra metadata to a log entry.
func With(key string, value any) LogOptions {
	return func(log *log) {
		log.addMetadata(key, value)
	}
}

// WithMessage adds a message to a log entry.
func WithMessage(msg string) LogOptions {
	return func(log *log) {
		log.addMetadata("message", msg)
	}
}

// isZero checks if the log entry has no level or data.
func (l *log) isZero() bool {
	return len(l.data) == 0 || l.level == ""
}

// addMetadata adds a key-value pair to the log's metadata.
func (l *log) addMetadata(k string, v any) {
	if strings.TrimSpace(k) != "" {
		l.data = append(l.data, logData{
			key:   k,
			value: v,
		})
	}
}

// simple generates a simple log message as a string.
func (l *log) simple(formatter TimeFormatter) string {
	var message strings.Builder

	// Add timestamp
	message.WriteString(formatter(l.timestamp, "[2006-01-02 15:04:05 Z0700] "))

	// Add level
	message.WriteString(fmt.Sprintf("%5s ", l.level))

	// Add metadata
	for _, item := range l.data {
		message.WriteString(fmt.Sprintf(`%s: "%v" `, item.key, item.value))
	}

	return message.String()
}

// structured generates a structured log message as a JSON-like byte slice.
func (l *log) structured(formatter TimeFormatter) []byte {
	var result []string

	// Helper function to add a key-value pair to the result.
	addItem := func(k string, v any) {
		encoded, _ := json.Marshal(v)
		if len(result) > 0 {
			result = append(result, fmt.Sprintf(` "%s": %s`, k, string(encoded)))
		} else {
			result = append(result, fmt.Sprintf(`"%s": %s`, k, string(encoded)))
		}
	}

	// Add level
	addItem("lvl", l.level)

	// Add timestamp
	addItem("ts", l.timestamp.Unix())
	addItem("dt", formatter(l.timestamp, "2006-01-02 15:04:05 Z0700"))

	// Add metadata
	for _, item := range l.data {
		addItem(item.key, item.value)
	}

	return []byte("{ " + strings.Join(result, ",") + " }")
}

// print outputs the log entry to the console with formatting.
func (l *log) print(formatter TimeFormatter) {
	// Print level with color formatting
	switch l.level {
	case levelDebug:
		console.PrintF("@Bp{%5s} ", l.level)
	case levelInfo:
		console.PrintF("@Bb{%5s} ", l.level)
	case levelWarn:
		console.PrintF("@By{%5s} ", l.level)
	case levelError, levelPanic:
		console.PrintF("@Br{%5s} ", l.level)
	}

	// Print timestamp
	console.PrintF("[@I{%s}] ", formatter(l.timestamp, "2006-01-02 15:04:05 Z0700"))

	// Print metadata
	for _, item := range l.data {
		encoded, _ := json.Marshal(item.value)
		console.PrintF(`@U{%s}: @Bg{%s} `, item.key, string(encoded))
	}

	console.PrintF("\n")
}
