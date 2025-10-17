package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestQuietLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := New(LevelQuiet, &buf)

	logger.Info("info message")
	logger.Debug("debug message")
	logger.Error("error message")

	output := buf.String()
	if strings.Contains(output, "info message") {
		t.Error("Quiet logger should not show info messages")
	}
	if strings.Contains(output, "debug message") {
		t.Error("Quiet logger should not show debug messages")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Quiet logger should show error messages")
	}
}

func TestNormalLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := New(LevelNormal, &buf)

	logger.Info("info message")
	logger.Debug("debug message")
	logger.Error("error message")

	output := buf.String()
	if !strings.Contains(output, "info message") {
		t.Error("Normal logger should show info messages")
	}
	if strings.Contains(output, "debug message") {
		t.Error("Normal logger should not show debug messages")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Normal logger should show error messages")
	}
}

func TestVerboseLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := New(LevelVerbose, &buf)

	logger.Info("info message")
	logger.Debug("debug message")
	logger.Error("error message")

	output := buf.String()
	if !strings.Contains(output, "info message") {
		t.Error("Verbose logger should show info messages")
	}
	if !strings.Contains(output, "debug message") {
		t.Error("Verbose logger should show debug messages")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Verbose logger should show error messages")
	}
}

func TestProgressIndicators(t *testing.T) {
	var buf bytes.Buffer
	logger := New(LevelNormal, &buf)

	logger.StartStep("Testing step")
	output := buf.String()
	if !strings.Contains(output, "⏳") {
		t.Error("StartStep should show hourglass emoji")
	}
	if !strings.Contains(output, "Testing step") {
		t.Error("StartStep should show step name")
	}

	buf.Reset()
	logger.CompleteStep("Testing step")
	output = buf.String()
	if !strings.Contains(output, "✓") {
		t.Error("CompleteStep should show checkmark")
	}

	buf.Reset()
	logger.FailStep("Testing step")
	output = buf.String()
	if !strings.Contains(output, "✗") {
		t.Error("FailStep should show X mark")
	}
}
