package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// ValidatePullSecret checks if the pull secret file exists and is valid JSON
func ValidatePullSecret(path string) error {
	if path == "" {
		return fmt.Errorf("pull secret path is empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read pull secret: %w", err)
	}

	var js interface{}
	if err := json.Unmarshal(data, &js); err != nil {
		return fmt.Errorf("pull secret is not valid JSON: %w", err)
	}

	return nil
}

// CheckPrerequisites validates that required tools are available
func CheckPrerequisites() error {
	// Check for oc command
	if _, err := exec.LookPath("oc"); err != nil {
		return fmt.Errorf("'oc' command not found in PATH. Please install OpenShift CLI")
	}

	return nil
}
