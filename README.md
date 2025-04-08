# Logger Package Documentation

![GitHub Tag](https://img.shields.io/github/v/tag/go-universal/logger?sort=semver&label=version)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-universal/logger.svg)](https://pkg.go.dev/github.com/go-universal/logger)
[![License](https://img.shields.io/badge/license-ISC-blue.svg)](https://github.com/go-universal/logger/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-universal/logger)](https://goreportcard.com/report/github.com/go-universal/logger)
![Contributors](https://img.shields.io/github/contributors/go-universal/logger)
![Issues](https://img.shields.io/github/issues/go-universal/logger)

The `logger` package provides a flexible and efficient logging system for Go applications. It supports multiple logging levels, structured and simple formats, and customizable time formatting.

## Installation

```bash
go get github.com/go-universal/logger
```

Here is an example of how to use the `logger` package:

```go
package main

import (
    "fmt"
    "github.com/go-universal/logger"
)

func main() {
    // Create a logger instance
    log, err := logger.NewLogger().
        SetBufferSize(100).
        Path("./logs").
        Prefix("app_").
        Extension("log").
        Daily().
        Simple().
        StdFormatter().
        Logger()

    if err != nil {
        fmt.Println("Failed to initialize logger:", err)
        return
    }

    // Log messages at different levels
    log.Debug(
        logger.With("Name", "John Doe"),
        logger.With("Age", 30),
        logger.WithMessage("Debugging application"),
    )

    log.Info(
        logger.With("Name", "Jane Doe"),
        logger.With("Age", 25),
        logger.WithMessage("Application started"),
    )

    log.Warn(
        logger.With("Name", "Jim Doe"),
        logger.With("Age", 40),
        logger.WithMessage("Potential issue detected"),
    )

    log.Error(
        logger.With("Name", "Jake Doe"),
        logger.With("Age", 50),
        logger.WithMessage("An error occurred"),
    )

    // Flush logs before exiting
    log.Sync()
}
```

## Types and Functions

### Logger

Defines methods for logging at various levels:

- **`Debug(options ...LogOptions)`**: Logs a debug-level message (development only).
- **`Info(options ...LogOptions)`**: Logs an info-level message.
- **`Warn(options ...LogOptions)`**: Logs a warning-level message.
- **`Error(options ...LogOptions)`**: Logs an error-level message.
- **`Panic(options ...LogOptions)`**: Logs a panic-level message.
- **`Sync()`**: Flushes any buffered log entries.

### Logger Builder

Used to configure and create a `Logger` instance:

- **`NewLogger()`**: Initializes a `LoggerBuilder` with default settings.
- **`SetBufferSize(size uint)`**: Configures the buffer size for the logger's channel.
- **`Development()`**: Enables development mode for the logger.
- **`Production()`**: Enables production mode for the logger.
- **`Simple()`**: Enables simple formatting for the logger.
- **`Structured()`**: Enables structured formatting for the logger.
- **`Silent()`**: Makes the logger silent (no output).
- **`Path(root string)`**: Sets the root directory for log files.
- **`Prefix(prefix string)`**: Sets the prefix for log file names.
- **`Extension(ext string)`**: Sets the extension for log file names.
- **`Daily()`**: Sets the log file layout to daily.
- **`Monthly()`**: Sets the log file layout to monthly.
- **`CustomLayout(layout string)`**: Sets a custom layout for log file names.
- **`StdFormatter()`**: Sets the logger to use the standard time formatter.
- **`JalaaliFormatter()`**: Sets the logger to use the Jalaali time formatter.
- **`CustomFormatter(formatter TimeFormatter)`**: Sets a custom time formatter for the logger.
- **`Logger()`**: Creates and returns a `Logger` instance.

### Log Options

Functions to modify log entries:

- **`With(key string, value any)`**: Adds extra metadata to a log entry.
- **`WithMessage(msg string)`**: Adds a message to a log entry.

### TimeFormatter

A function signature for time formatting:

```go
type TimeFormatter func(t time.Time, layout string) string
```

Built-in formatters:

- **`StdFormatter`**: Formats time using the standard library.
- **`JalaaliFormatter`**: Formats time using the Jalaali calendar.

### Log Levels

Defines the severity levels for log messages:

- `DEBUG`, `INFO`, `WARN`, `ERROR`, `PANIC`.

## License

This project is licensed under the ISC License. See the [LICENSE](LICENSE) file for details.
