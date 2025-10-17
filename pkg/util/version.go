package util

import (
	"fmt"
	"strings"
)

// ExtractVersionArch extracts the version-arch portion from a release image URL
// Example: "quay.io/openshift-release-dev/ocp-release:4.12.0-x86_64" -> "4.12.0-x86_64"
func ExtractVersionArch(releaseImage string) (string, error) {
	if releaseImage == "" {
		return "", fmt.Errorf("release image cannot be empty")
	}

	// Find the last colon which separates the tag
	parts := strings.Split(releaseImage, ":")
	if len(parts) < 2 {
		return "", fmt.Errorf("release image must contain a tag (e.g., :4.12.0-x86_64)")
	}

	// The tag is the last part
	tag := parts[len(parts)-1]
	if tag == "" {
		return "", fmt.Errorf("release image tag cannot be empty")
	}

	return tag, nil
}
