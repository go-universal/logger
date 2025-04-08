package logger

import (
	"time"

	"github.com/go-universal/jalaali"
)

// TimeFormatter function signature for time formatter
type TimeFormatter func(t time.Time, layout string) string

// StdFormatter standard time formatter instance.
func StdFormatter(ts time.Time, layout string) string {
	return ts.Format(layout)
}

// JalaaliFormatter jalaali time formatter instance.
func JalaaliFormatter(ts time.Time, layout string) string {
	return jalaali.New(ts).Format(layout)
}
