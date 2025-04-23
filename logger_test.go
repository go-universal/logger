package logger_test

import (
	"fmt"
	"testing"

	"github.com/go-universal/logger"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	l, err := logger.NewLogger().
		SetBufferSize(10).
		Extension("log").
		Daily().
		SetSimple(true).
		JalaaliFormatter().
		Logger()

	require.NoError(t, err, "Logger initialization should not fail")

	for i := 0; i < 150; i++ {
		switch i % 5 {
		case 0:
			l.Debug(
				logger.With("Name", fmt.Sprintf("John Doe %d", i)),
				logger.With("Age", i+20),
				logger.WithMessage(fmt.Sprintf("Something happened %d", i)),
			)
		case 1:
			l.Info(
				logger.With("Name", fmt.Sprintf("Jane Doe %d", i)),
				logger.With("Age", i+30),
				logger.WithMessage(fmt.Sprintf("Something interesting %d", i)),
			)
		case 2:
			l.Warn(
				logger.With("Name", fmt.Sprintf("Jim Doe %d", i)),
				logger.With("Age", i+40),
				logger.WithMessage(fmt.Sprintf("Something unusual %d", i)),
			)
		case 3:
			l.Error(
				logger.With("Name", fmt.Sprintf("Jake Doe %d", i)),
				logger.With("Age", i+50),
				logger.WithMessage(fmt.Sprintf("Something bad %d", i)),
			)
		case 4:
			l.Panic(
				logger.With("Name", fmt.Sprintf("Jake Doe %d", i)),
				logger.With("Age", i+50),
				logger.WithMessage(fmt.Sprintf("Something bad %d", i)),
			)
		}
	}

	l.Sync()
}
