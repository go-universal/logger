package logger

import (
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Logger defines methods for logging at various levels.
// It supports different modes for production and development environments.
type Logger interface {
	// Debug logs a debug-level message (development only).
	Debug(options ...LogOptions)

	// Info logs an info-level message.
	Info(options ...LogOptions)

	// Warn logs a warning-level message.
	Warn(options ...LogOptions)

	// Error logs an error-level message.
	Error(options ...LogOptions)

	// Panic logs a panic-level message.
	Panic(options ...LogOptions)

	// Sync flushes any buffered log entries.
	// Call this before exiting to ensure all logs are written.
	Sync()
}

type logger struct {
	channel chan log
	flushed atomic.Bool
	wg      sync.WaitGroup

	path string
	file *os.File

	dev       bool
	simple    bool
	silent    bool
	root      string
	prefix    string
	ext       string
	layout    string
	formatter TimeFormatter
}

func (l *logger) Debug(options ...LogOptions) {
	if l.flushed.Load() || !l.dev {
		return
	}

	l.queueLog(levelDebug, options...)
}

func (l *logger) Info(options ...LogOptions) {
	if l.flushed.Load() {
		return
	}

	l.queueLog(levelInfo, options...)
}

func (l *logger) Warn(options ...LogOptions) {
	if l.flushed.Load() {
		return
	}

	l.queueLog(levelWarn, options...)
}

func (l *logger) Error(options ...LogOptions) {
	if l.flushed.Load() {
		return
	}

	l.queueLog(levelError, options...)
}

func (l *logger) Panic(options ...LogOptions) {
	if l.flushed.Load() {
		return
	}

	l.queueLog(levelPanic, options...)
}

func (l *logger) Sync() {
	if l.flushed.CompareAndSwap(false, true) {
		close(l.channel)
		l.wg.Wait()
	}
}

// run processes log entries from the channel and writes them to the appropriate destination.
func (l *logger) run() {
	defer l.closeFile()
	defer l.wg.Done()

	for log := range l.channel {
		// Generate file path and skip if empty
		filePath := l.filePath(log.timestamp)
		if filePath == "" {
			continue
		}

		// Generate log message and skip if empty
		var message string
		if l.simple {
			message = log.simple(l.formatter)
		} else {
			message = string(log.structured(l.formatter))
		}
		if message == "" {
			continue
		}

		// Print to console if in development mode and not silent
		if l.dev && !l.silent {
			log.print(l.formatter)
		}

		// Write to file
		l.openFile(filePath)
		if l.file != nil {
			if _, err := l.file.WriteString(message + "\n"); err != nil {
				// Reopen file on error to prevent being stuck
				l.file.Close()
				l.file = nil
				l.path = ""
			}
		}
	}
}

// queueLog creates a new log entry and queues it for processing.
func (l *logger) queueLog(level Level, options ...LogOptions) {
	log := l.newLog(level, options...)
	if log.isZero() {
		return
	}

	l.channel <- log
}

// newLog creates a new log entry with the specified level and options.
func (l *logger) newLog(level Level, opts ...LogOptions) log {
	log := &log{
		level:     level,
		timestamp: time.Now(),
		data:      make([]logData, 0),
	}

	for _, opt := range opts {
		opt(log)
	}

	return *log
}

// filePath generates the file path for the log file based on the timestamp.
func (l *logger) filePath(ts time.Time) string {
	// Generate file name
	name := l.prefix + l.formatter(ts, l.layout)
	if name == "" {
		return ""
	}

	// Append file extension
	if ext := strings.TrimLeft(l.ext, "."); ext != "" {
		name += "." + ext
	}

	// Return full path
	return path.Join(l.root, name)
}

// openFile opens the log file for writing. Reopens if the path changes.
func (l *logger) openFile(filePath string) {
	if filePath == "" || (l.path == filePath && l.file != nil) {
		return
	}

	// Close the old file if the path has changed
	if l.path != filePath && l.file != nil {
		l.file.Close()
		l.file = nil
	}

	// Open the new file
	if f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		l.path = ""
		l.file = nil
	} else {
		l.path = filePath
		l.file = f
	}
}

// closeFile closes the currently open log file.
func (l *logger) closeFile() {
	if l.file != nil {
		l.file.Close()
	}
}
