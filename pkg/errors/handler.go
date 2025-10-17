package errors

import (
	"fmt"
	"strings"
)

type StepError struct {
	StepName string
	Error    error
}

type Summary struct {
	Successful []string
	Failed     []StepError
}

func NewSummary() *Summary {
	return &Summary{
		Successful: []string{},
		Failed:     []StepError{},
	}
}

func (s *Summary) AddSuccess(stepName string) {
	s.Successful = append(s.Successful, stepName)
}

func (s *Summary) AddError(stepName string, err error) {
	s.Failed = append(s.Failed, StepError{
		StepName: stepName,
		Error:    err,
	})
}

func (s *Summary) HasErrors() bool {
	return len(s.Failed) > 0
}

func (s *Summary) String() string {
	var sb strings.Builder

	sb.WriteString("\n=== Installation Summary ===\n\n")

	if len(s.Successful) > 0 {
		sb.WriteString("✓ Successful steps:\n")
		for _, step := range s.Successful {
			sb.WriteString(fmt.Sprintf("  - %s\n", step))
		}
		sb.WriteString("\n")
	}

	if len(s.Failed) > 0 {
		sb.WriteString("✗ Failed steps:\n")
		for _, stepErr := range s.Failed {
			sb.WriteString(fmt.Sprintf("  - %s: %v\n", stepErr.StepName, stepErr.Error))
		}
		sb.WriteString("\n")
	}

	if s.HasErrors() {
		sb.WriteString("Overall status: PARTIAL SUCCESS (some steps failed)\n")
	} else if len(s.Successful) > 0 {
		sb.WriteString("Overall status: SUCCESS\n")
	} else {
		sb.WriteString("Overall status: NO STEPS EXECUTED\n")
	}

	return sb.String()
}
