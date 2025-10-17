package steps

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/clobrano/ccoctl-sso/pkg/config"
	"github.com/clobrano/ccoctl-sso/pkg/logger"
	"github.com/clobrano/ccoctl-sso/pkg/util"
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
		s.cfg.ReleaseImage,
	}

	return util.RunCommand(s.executor, "oc", args...)
}

// Step2ExtractBinaries extracts openshift-install and ccoctl binaries
type Step2ExtractBinaries struct {
	*BaseStep
}

func NewStep2(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step2ExtractBinaries, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step2ExtractBinaries{BaseStep: base}, nil
}

func (s *Step2ExtractBinaries) Name() string {
	return "Extract binaries"
}

func (s *Step2ExtractBinaries) Execute() error {
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
		s.cfg.ReleaseImage,
	}
	if err := util.RunCommand(s.executor, "oc", args...); err != nil {
		return fmt.Errorf("failed to extract openshift-install: %w", err)
	}

	// Make it executable
	os.Chmod(installBinPath, 0755)

	// Extract ccoctl
	ccoctlPath := filepath.Join(binPath, "ccoctl")

	// Get CCO image
	ccoImageArgs := []string{"adm", "release", "info", "--image-for=cloud-credential-operator", s.cfg.ReleaseImage}
	ccoImage, err := s.executor.Execute("oc", ccoImageArgs...)
	if err != nil {
		return fmt.Errorf("failed to get CCO image: %w", err)
	}

	// Extract ccoctl from CCO image
	extractArgs := []string{
		"image", "extract",
		ccoImage,
		"--file=/usr/bin/ccoctl:" + ccoctlPath,
	}
	if err := util.RunCommand(s.executor, "oc", extractArgs...); err != nil {
		return fmt.Errorf("failed to extract ccoctl: %w", err)
	}

	// Make it executable
	os.Chmod(ccoctlPath, 0755)

	return nil
}

// Step3CreateConfig runs openshift-install create install-config
type Step3CreateConfig struct {
	*BaseStep
}

func NewStep3(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step3CreateConfig, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step3CreateConfig{BaseStep: base}, nil
}

func (s *Step3CreateConfig) Name() string {
	return "Create install-config.yaml"
}

func (s *Step3CreateConfig) Execute() error {
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

// Step4SetCredentialsMode appends credentialsMode: Manual to install-config.yaml
type Step4SetCredentialsMode struct {
	*BaseStep
}

func NewStep4(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step4SetCredentialsMode, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step4SetCredentialsMode{BaseStep: base}, nil
}

func (s *Step4SetCredentialsMode) Name() string {
	return "Set credentialsMode to Manual"
}

func (s *Step4SetCredentialsMode) Execute() error {
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

// Step5CreateManifests runs openshift-install create manifests
type Step5CreateManifests struct {
	*BaseStep
}

func NewStep5(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step5CreateManifests, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step5CreateManifests{BaseStep: base}, nil
}

func (s *Step5CreateManifests) Name() string {
	return "Create manifests"
}

func (s *Step5CreateManifests) Execute() error {
	versionDir := filepath.Join("artifacts", s.versionArch)
	installBin := util.GetBinaryPath(s.versionArch, "openshift-install")
	args := []string{"create", "manifests", "--dir", versionDir}

	return util.RunCommand(s.executor, installBin, args...)
}

// Additional steps will follow the same pattern...
// For brevity, I'll implement the remaining steps in separate files
