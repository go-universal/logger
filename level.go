package logger

// Level represents the severity level of a log message.
type Level string

// Log levels.
const (
	levelDebug Level = "DEBUG"
	levelInfo  Level = "INFO"
	levelWarn  Level = "WARN"
	levelError Level = "ERROR"
	levelPanic Level = "PANIC"
)
