package util

import "testing"

func TestExtractVersionArch(t *testing.T) {
	tests := []struct {
		name          string
		releaseImage  string
		expected      string
		shouldSucceed bool
	}{
		{
			name:          "standard release image",
			releaseImage:  "quay.io/openshift-release-dev/ocp-release:4.12.0-x86_64",
			expected:      "4.12.0-x86_64",
			shouldSucceed: true,
		},
		{
			name:          "release with fc version",
			releaseImage:  "quay.io/openshift-release-dev/ocp-release:4.10.0-fc.4-x86_64",
			expected:      "4.10.0-fc.4-x86_64",
			shouldSucceed: true,
		},
		{
			name:          "aarch64 architecture",
			releaseImage:  "quay.io/openshift-release-dev/ocp-release:4.13.1-aarch64",
			expected:      "4.13.1-aarch64",
			shouldSucceed: true,
		},
		{
			name:          "no tag",
			releaseImage:  "quay.io/openshift-release-dev/ocp-release",
			expected:      "",
			shouldSucceed: false,
		},
		{
			name:          "empty string",
			releaseImage:  "",
			expected:      "",
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractVersionArch(tt.releaseImage)
			if tt.shouldSucceed {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %q but got %q", tt.expected, result)
				}
			} else {
				if err == nil {
					t.Error("Expected error but got success")
				}
			}
		})
	}
}
