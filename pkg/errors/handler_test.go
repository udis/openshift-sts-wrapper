package errors

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorSummary(t *testing.T) {
	summary := NewSummary()

	summary.AddSuccess("Step 1")
	summary.AddSuccess("Step 2")
	summary.AddError("Step 3", errors.New("something failed"))
	summary.AddSuccess("Step 4")

	if len(summary.Successful) != 3 {
		t.Errorf("Expected 3 successful steps, got %d", len(summary.Successful))
	}

	if len(summary.Failed) != 1 {
		t.Errorf("Expected 1 failed step, got %d", len(summary.Failed))
	}

	if summary.Failed[0].StepName != "Step 3" {
		t.Error("Failed step should be 'Step 3'")
	}

	if summary.HasErrors() != true {
		t.Error("Summary should have errors")
	}
}

func TestErrorSummaryString(t *testing.T) {
	summary := NewSummary()
	summary.AddSuccess("Extract credentials")
	summary.AddError("Extract binaries", errors.New("download failed"))
	summary.AddSuccess("Create config")

	output := summary.String()

	if !strings.Contains(output, "Extract credentials") {
		t.Error("Summary should contain successful step")
	}
	if !strings.Contains(output, "Extract binaries") {
		t.Error("Summary should contain failed step")
	}
	if !strings.Contains(output, "download failed") {
		t.Error("Summary should contain error message")
	}
}

func TestEmptySummary(t *testing.T) {
	summary := NewSummary()

	if summary.HasErrors() {
		t.Error("Empty summary should not have errors")
	}

	if len(summary.Successful) != 0 {
		t.Error("Empty summary should have no successful steps")
	}
}
