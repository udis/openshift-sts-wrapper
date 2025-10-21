package util

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// InstallConfig represents the minimal structure we need from install-config.yaml
type InstallConfig struct {
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Platform struct {
		AWS struct {
			Region string `yaml:"region"`
		} `yaml:"aws"`
	} `yaml:"platform"`
}

// ReadInstallConfig reads and parses install-config.yaml
func ReadInstallConfig(path string) (*InstallConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read install-config.yaml: %w", err)
	}

	var config InstallConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse install-config.yaml: %w", err)
	}

	return &config, nil
}

// ExtractClusterNameAndRegion reads install-config.yaml and returns the cluster name and region
func ExtractClusterNameAndRegion(installConfigPath string) (clusterName string, region string, err error) {
	config, err := ReadInstallConfig(installConfigPath)
	if err != nil {
		return "", "", err
	}

	if config.Metadata.Name == "" {
		return "", "", fmt.Errorf("cluster name not found in install-config.yaml")
	}

	if config.Platform.AWS.Region == "" {
		return "", "", fmt.Errorf("AWS region not found in install-config.yaml")
	}

	return config.Metadata.Name, config.Platform.AWS.Region, nil
}
