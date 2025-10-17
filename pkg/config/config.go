package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ReleaseImage   string `yaml:"releaseImage"`
	ClusterName    string `yaml:"clusterName"`
	AwsRegion      string `yaml:"awsRegion"`
	PullSecretPath string `yaml:"pullSecretPath"`
	PrivateBucket  bool   `yaml:"privateBucket"`
	OutputDir      string `yaml:"outputDir"`
	StartFromStep  int    `yaml:"startFromStep"`
}

// LoadFromFile loads configuration from a YAML file
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() *Config {
	return &Config{
		ReleaseImage:   os.Getenv("OPENSHIFT_STS_RELEASE_IMAGE"),
		ClusterName:    os.Getenv("OPENSHIFT_STS_CLUSTER_NAME"),
		AwsRegion:      os.Getenv("OPENSHIFT_STS_AWS_REGION"),
		PullSecretPath: os.Getenv("OPENSHIFT_STS_PULL_SECRET_PATH"),
		PrivateBucket:  os.Getenv("OPENSHIFT_STS_PRIVATE_BUCKET") == "true",
		OutputDir:      os.Getenv("OPENSHIFT_STS_OUTPUT_DIR"),
	}
}

// Merge merges another config into this one, with the other config taking precedence
func (c *Config) Merge(other *Config) {
	if other.ReleaseImage != "" {
		c.ReleaseImage = other.ReleaseImage
	}
	if other.ClusterName != "" {
		c.ClusterName = other.ClusterName
	}
	if other.AwsRegion != "" {
		c.AwsRegion = other.AwsRegion
	}
	if other.PullSecretPath != "" {
		c.PullSecretPath = other.PullSecretPath
	}
	if other.PrivateBucket {
		c.PrivateBucket = other.PrivateBucket
	}
	if other.OutputDir != "" {
		c.OutputDir = other.OutputDir
	}
	if other.StartFromStep > 0 {
		c.StartFromStep = other.StartFromStep
	}
}

// ValidateConfig validates that required fields are set
func ValidateConfig(cfg *Config) error {
	if cfg.ReleaseImage == "" {
		return fmt.Errorf("release image is required")
	}
	if cfg.ClusterName == "" {
		return fmt.Errorf("cluster name is required")
	}
	if cfg.AwsRegion == "" {
		return fmt.Errorf("AWS region is required")
	}
	return nil
}

// SetDefaults sets default values for optional fields
func (c *Config) SetDefaults() {
	if c.OutputDir == "" {
		c.OutputDir = "_output"
	}
	if c.PullSecretPath == "" {
		c.PullSecretPath = "pull-secret.json"
	}
}
