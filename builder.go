package logger

import (
	"os"
	"strings"
)

// LoggerBuilder is used to configure and create a Logger instance.
type LoggerBuilder struct {
	buf       uint          // Buffer size for the logger's channel.
	dev       bool          // Indicates if the logger is in development mode.
	simple    bool          // Indicates if the logger uses a simple format.
	silent    bool          // Indicates if the logger is silent.
	root      string        // Root directory for log files.
	prefix    string        // Prefix for log file names.
	ext       string        // Extension for log file names.
	layout    string        // Layout for log file names.
	formatter TimeFormatter // Formatter for log timestamps.
}

// NewLogger initializes a LoggerBuilder with default settings.
func NewLogger() *LoggerBuilder {
	return &LoggerBuilder{
		buf:       100,
		dev:       true,
		simple:    false,
		silent:    false,
		root:      "./logs",
		prefix:    "",
		ext:       "",
		layout:    "2006-01-02",
		formatter: StdFormatter,
	}
}

// SetBufferSize configures the buffer size for the logger's channel.
func (b *LoggerBuilder) SetBufferSize(size uint) *LoggerBuilder {
	b.buf = size
	return b
}

// SetEnv set development or production mode
func (b *LoggerBuilder) SetEnv(dev bool) *LoggerBuilder {
	b.dev = dev
	return b
}

// SetSimple set simple formatting for the logger.
func (b *LoggerBuilder) SetSimple(simple bool) *LoggerBuilder {
	b.simple = simple
	return b
}

// SetSilent set silent mode for logger (no print to console).
func (b *LoggerBuilder) SetSilent(silent bool) *LoggerBuilder {
	b.silent = silent
	return b
}

// Path sets the root directory for log files. Ignores empty input.
func (b *LoggerBuilder) Path(root string) *LoggerBuilder {
	root = strings.TrimSpace(root)
	if root != "" {
		b.root = root
	}
	return b
}

// Prefix sets the prefix for log file names. Ignores empty input.
func (b *LoggerBuilder) Prefix(prefix string) *LoggerBuilder {
	if strings.TrimSpace(prefix) != "" {
		b.prefix = prefix
	}
	return b
}

// Extension sets the extension for log file names. Ignores empty input.
func (b *LoggerBuilder) Extension(ext string) *LoggerBuilder {
	ext = strings.TrimSpace(ext)
	if ext != "" {
		b.ext = ext
	}
	return b
}

// Daily sets the log file layout to daily.
func (b *LoggerBuilder) Daily() *LoggerBuilder {
	b.layout = "2006-01-02"
	return b
}

// Monthly sets the log file layout to monthly.
func (b *LoggerBuilder) Monthly() *LoggerBuilder {
	b.layout = "2006-01"
	return b
}

// CustomLayout sets a custom layout for log file names. Ignores empty input.
func (b *LoggerBuilder) CustomLayout(layout string) *LoggerBuilder {
	layout = strings.TrimSpace(layout)
	if layout != "" {
		b.layout = layout
	}
	return b
}

// StdFormatter sets the logger to use the standard time formatter.
func (b *LoggerBuilder) StdFormatter() *LoggerBuilder {
	b.formatter = StdFormatter
	return b
}

// JalaaliFormatter sets the logger to use the Jalaali time formatter.
func (b *LoggerBuilder) JalaaliFormatter() *LoggerBuilder {
	b.formatter = JalaaliFormatter
	return b
}

// CustomFormatter sets a custom time formatter for the logger. Ignores nil input.
func (b *LoggerBuilder) CustomFormatter(formatter TimeFormatter) *LoggerBuilder {
	if formatter != nil {
		b.formatter = formatter
	}
	return b
}

// Logger creates and returns a Logger instance based on the builder's configuration.
func (b *LoggerBuilder) Logger() (Logger, error) {
	// Ensure the log directory exists.
	err := os.MkdirAll(b.root, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	// Create and initialize the logger.
	logger := &logger{
		channel:   make(chan log, b.buf),
		dev:       b.dev,
		simple:    b.simple,
		silent:    b.silent,
		root:      b.root,
		prefix:    b.prefix,
		ext:       b.ext,
		layout:    b.layout,
		formatter: b.formatter,
	}
	logger.wg.Add(1)
	go logger.run()
	return logger, nil
}
