package steps

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.cee.redhat.com/clobrano/ccoctl-sso/pkg/config"
	"gitlab.cee.redhat.com/clobrano/ccoctl-sso/pkg/logger"
	"gitlab.cee.redhat.com/clobrano/ccoctl-sso/pkg/util"
)

// Step represents a single installation step
type Step interface {
	Name() string
	Execute() error
}

// BaseStep contains common fields for all steps
type BaseStep struct {
	cfg         *config.Config
	log         *logger.Logger
	executor    util.CommandExecutor
	versionArch string
}

func newBaseStep(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*BaseStep, error) {
	versionArch, err := util.ExtractVersionArch(cfg.ReleaseImage)
	if err != nil {
		return nil, err
	}

	return &BaseStep{
		cfg:         cfg,
		log:         log,
		executor:    executor,
		versionArch: versionArch,
	}, nil
}

// Step1ExtractCredReqs extracts credentials requests from the release image
type Step1ExtractCredReqs struct {
	*BaseStep
}

func NewStep1(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step1ExtractCredReqs, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step1ExtractCredReqs{BaseStep: base}, nil
}

func (s *Step1ExtractCredReqs) Name() string {
	return "Extract credentials requests"
}

func (s *Step1ExtractCredReqs) Execute() error {
	credreqsPath := util.GetCredReqsPath(s.versionArch)
	if err := util.EnsureDir(credreqsPath); err != nil {
		return fmt.Errorf("failed to create credreqs directory: %w", err)
	}

	args := []string{
		"adm", "release", "extract",
		"--credentials-requests",
		"--cloud=aws",
		"--to=" + credreqsPath,
		"--registry-config=" + s.cfg.PullSecretPath,
		s.cfg.ReleaseImage,
	}

	return util.RunCommand(s.executor, "oc", args...)
}

// Step2ExtractOpenshiftInstall extracts openshift-install binary
type Step2ExtractOpenshiftInstall struct {
	*BaseStep
}

func NewStep2(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step2ExtractOpenshiftInstall, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step2ExtractOpenshiftInstall{BaseStep: base}, nil
}

func (s *Step2ExtractOpenshiftInstall) Name() string {
	return "Extract openshift-install binary"
}

func (s *Step2ExtractOpenshiftInstall) Execute() error {
	binPath := filepath.Join("artifacts", s.versionArch, "bin")
	if err := util.EnsureDir(binPath); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Extract openshift-install
	installBinPath := filepath.Join(binPath, "openshift-install")
	args := []string{
		"adm", "release", "extract",
		"--command=openshift-install",
		"--to=" + binPath,
		"--registry-config=" + s.cfg.PullSecretPath,
		s.cfg.ReleaseImage,
	}
	if err := util.RunCommand(s.executor, "oc", args...); err != nil {
		return fmt.Errorf("failed to extract openshift-install: %w", err)
	}

	// Make it executable
	os.Chmod(installBinPath, 0755)

	return nil
}

// Step3ExtractCcoctl extracts ccoctl binary
type Step3ExtractCcoctl struct {
	*BaseStep
}

func NewStep3(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step3ExtractCcoctl, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step3ExtractCcoctl{BaseStep: base}, nil
}

func (s *Step3ExtractCcoctl) Name() string {
	return "Extract ccoctl binary"
}

func (s *Step3ExtractCcoctl) Execute() error {
	binPath := filepath.Join("artifacts", s.versionArch, "bin")
	ccoctlPath := filepath.Join(binPath, "ccoctl")

	// Get CCO image
	ccoImageArgs := []string{"adm", "release", "info", "--image-for=cloud-credential-operator", "--registry-config=" + s.cfg.PullSecretPath, s.cfg.ReleaseImage}
	ccoImage, err := s.executor.Execute("oc", ccoImageArgs...)
	if err != nil {
		return fmt.Errorf("failed to get CCO image: %w", err)
	}

	// Trim whitespace from CCO image reference
	ccoImage = strings.TrimSpace(ccoImage)

	// Extract ccoctl from CCO image
	extractArgs := []string{
		"image", "extract",
		ccoImage,
		"--file=/usr/bin/ccoctl:",
		"--registry-config=" + s.cfg.PullSecretPath,
	}
	if err := util.RunCommand(s.executor, "oc", extractArgs...); err != nil {
		return fmt.Errorf("failed to extract ccoctl: %w", err)
	}

	// Make it executable
	os.Chmod(ccoctlPath, 0755)

	return nil
}

// Step4CreateConfig runs openshift-install create install-config
type Step4CreateConfig struct {
	*BaseStep
}

func NewStep4(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step4CreateConfig, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step4CreateConfig{BaseStep: base}, nil
}

func (s *Step4CreateConfig) Name() string {
	return "Create install-config.yaml"
}

func (s *Step4CreateConfig) Execute() error {
	// Ensure the version-specific directory exists
	versionDir := filepath.Join("artifacts", s.versionArch)
	if err := util.EnsureDir(versionDir); err != nil {
		return err
	}

	// Run openshift-install create install-config
	// Note: This is interactive in real usage, but mocked in tests
	installBin := util.GetBinaryPath(s.versionArch, "openshift-install")
	args := []string{"create", "install-config", "--dir", versionDir}

	return util.RunCommand(s.executor, installBin, args...)
}

// Step5SetCredentialsMode appends credentialsMode: Manual to install-config.yaml
type Step5SetCredentialsMode struct {
	*BaseStep
}

func NewStep5(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step5SetCredentialsMode, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step5SetCredentialsMode{BaseStep: base}, nil
}

func (s *Step5SetCredentialsMode) Name() string {
	return "Set credentialsMode to Manual"
}

func (s *Step5SetCredentialsMode) Execute() error {
	configPath := util.GetInstallConfigPath(s.versionArch)

	// Read existing config
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read install-config.yaml: %w", err)
	}

	// Append credentialsMode if not present
	configStr := string(content)
	if !util.FileContains(configPath, "credentialsMode:") {
		configStr += "\ncredentialsMode: Manual\n"
		if err := os.WriteFile(configPath, []byte(configStr), 0644); err != nil {
			return fmt.Errorf("failed to write install-config.yaml: %w", err)
		}
	}

	return nil
}

// Step6CreateManifests runs openshift-install create manifests
type Step6CreateManifests struct {
	*BaseStep
}

func NewStep6(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step6CreateManifests, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step6CreateManifests{BaseStep: base}, nil
}

func (s *Step6CreateManifests) Name() string {
	return "Create manifests"
}

func (s *Step6CreateManifests) Execute() error {
	versionDir := filepath.Join("artifacts", s.versionArch)
	installBin := util.GetBinaryPath(s.versionArch, "openshift-install")
	args := []string{"create", "manifests", "--dir", versionDir}

	return util.RunCommand(s.executor, installBin, args...)
}

// Additional steps will follow the same pattern...
// For brevity, I'll implement the remaining steps in separate files
