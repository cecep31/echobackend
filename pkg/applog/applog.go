package applog

import (
	"log/slog"
	"os"
	"time"
)

// Logger emits structured logs for a fixed component name.
// It resolves slog.Default() on each call, so package-level vars remain safe
// even when they are initialized before Setup().
type Logger struct {
	component string
}

// Setup installs the application-wide default slog handler on stdout.
// Level is Info by default, Debug when debug is true.
func Setup(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.RFC3339))
				}
			}
			return a
		},
	})
	slog.SetDefault(slog.New(handler))
}

// SetupFromEnv configures logging before full config load (APP_DEBUG or DEBUG).
func SetupFromEnv() {
	debug := os.Getenv("APP_DEBUG") == "true" || os.Getenv("DEBUG") == "true"
	Setup(debug)
}

// Default returns the configured application logger.
func Default() *slog.Logger {
	return slog.Default()
}

// Component returns a lazy component logger safe for package-level vars.
func Component(name string) Logger {
	return Logger{component: name}
}

func (l Logger) slog() *slog.Logger {
	return slog.Default().With("component", l.component)
}

// Slog returns the underlying slog logger for integrations such as GORM.
func (l Logger) Slog() *slog.Logger {
	return l.slog()
}

func (l Logger) Debug(msg string, args ...any) {
	l.slog().Debug(msg, args...)
}

func (l Logger) Info(msg string, args ...any) {
	l.slog().Info(msg, args...)
}

func (l Logger) Warn(msg string, args ...any) {
	l.slog().Warn(msg, args...)
}

func (l Logger) Error(msg string, args ...any) {
	l.slog().Error(msg, args...)
}
