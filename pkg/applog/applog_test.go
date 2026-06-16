package applog

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestComponentLoggerFormat(t *testing.T) {
	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))

	log := Component("database")
	log.Info("pool ready", "max_open", 30, "max_idle", 2)

	out := buf.String()
	if strings.Contains(out, `msg="INFO pool ready`) {
		t.Fatalf("unexpected duplicated level in message: %q", out)
	}
	if !strings.Contains(out, `msg="pool ready"`) {
		t.Fatalf("expected clean message field, got: %q", out)
	}
	if !strings.Contains(out, `component=database`) {
		t.Fatalf("expected component attribute, got: %q", out)
	}
}

func TestComponentLoggerBeforeSetup(t *testing.T) {
	earlyLog := Component("database")

	var buf bytes.Buffer
	Setup(false)
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))

	earlyLog.Info("pool ready", "max_open", 30)

	out := buf.String()
	if strings.Contains(out, `msg="INFO pool ready`) {
		t.Fatalf("duplicated INFO in message: %q", out)
	}
	if !strings.Contains(out, `msg="pool ready"`) {
		t.Fatalf("expected clean message field, got: %q", out)
	}
}
