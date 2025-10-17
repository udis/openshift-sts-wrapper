package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePullSecret(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		shouldError bool
	}{
		{
			name:        "valid JSON",
			content:     `{"auths":{"cloud.openshift.com":{"auth":"token"}}}`,
			shouldError: false,
		},
		{
			name:        "invalid JSON",
			content:     `{invalid json}`,
			shouldError: true,
		},
		{
			name:        "empty file",
			content:     ``,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".json")
			os.WriteFile(path, []byte(tt.content), 0644)

			err := ValidatePullSecret(path)
			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidatePullSecretNonExistent(t *testing.T) {
	err := ValidatePullSecret("/nonexistent/path/pull-secret.json")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestValidatePullSecretEmptyPath(t *testing.T) {
	err := ValidatePullSecret("")
	if err == nil {
		t.Error("Expected error for empty path")
	}
}
