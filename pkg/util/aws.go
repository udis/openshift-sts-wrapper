package util

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AWSCredentials holds AWS credentials from the credentials file
type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

// ReadAWSCredentials reads AWS credentials from ~/.aws/credentials for a given profile
func ReadAWSCredentials(profile string) (*AWSCredentials, error) {
	if profile == "" {
		profile = "default"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	credentialsPath := filepath.Join(homeDir, ".aws", "credentials")
	file, err := os.Open(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open credentials file: %w", err)
	}
	defer file.Close()

	creds := &AWSCredentials{}
	inTargetSection := false
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Check for section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.Trim(line, "[]")
			inTargetSection = (sectionName == profile)
			continue
		}

		// If we're in the target section, read the credentials
		if inTargetSection {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "aws_access_key_id":
				creds.AccessKeyID = value
			case "aws_secret_access_key":
				creds.SecretAccessKey = value
			case "aws_session_token":
				creds.SessionToken = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading credentials file: %w", err)
	}

	// Validate that we found at least the required credentials
	if creds.AccessKeyID == "" || creds.SecretAccessKey == "" {
		return nil, fmt.Errorf("profile '%s' not found or missing required credentials", profile)
	}

	return creds, nil
}

// GetAWSEnvVars returns environment variables for AWS credentials
func GetAWSEnvVars(profile string) ([]string, error) {
	creds, err := ReadAWSCredentials(profile)
	if err != nil {
		return nil, err
	}

	envVars := []string{
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", creds.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", creds.SecretAccessKey),
	}

	if creds.SessionToken != "" {
		envVars = append(envVars, fmt.Sprintf("AWS_SESSION_TOKEN=%s", creds.SessionToken))
	}

	return envVars, nil
}
