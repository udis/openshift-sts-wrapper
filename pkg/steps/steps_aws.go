package steps

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gitlab.cee.redhat.com/clobrano/ccoctl-sso/pkg/config"
	"gitlab.cee.redhat.com/clobrano/ccoctl-sso/pkg/logger"
	"gitlab.cee.redhat.com/clobrano/ccoctl-sso/pkg/util"
)

// Step6CreateAWSResources runs ccoctl to create AWS resources
type Step6CreateAWSResources struct {
	*BaseStep
}

func NewStep6(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step6CreateAWSResources, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step6CreateAWSResources{BaseStep: base}, nil
}

func (s *Step6CreateAWSResources) Name() string {
	return "Create AWS resources"
}

func (s *Step6CreateAWSResources) Execute() error {
	ccoctlBin := util.GetBinaryPath(s.versionArch, "ccoctl")
	credreqsPath := util.GetCredReqsPath(s.versionArch)

	args := []string{
		"aws", "create-all",
		"--name", s.cfg.ClusterName,
		"--region", s.cfg.AwsRegion,
		"--credentials-requests-dir", credreqsPath,
		"--output-dir", s.cfg.OutputDir,
	}

	if s.cfg.PrivateBucket {
		args = append(args, "--create-private-s3-bucket")
	}

	return util.RunCommand(s.executor, ccoctlBin, args...)
}

// Step7CopyManifests copies manifests from _output to manifests/
type Step7CopyManifests struct {
	*BaseStep
}

func NewStep7(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step7CopyManifests, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step7CopyManifests{BaseStep: base}, nil
}

func (s *Step7CopyManifests) Name() string {
	return "Copy manifests"
}

func (s *Step7CopyManifests) Execute() error {
	srcDir := filepath.Join(s.cfg.OutputDir, "manifests")
	dstDir := "manifests"

	if err := util.EnsureDir(dstDir); err != nil {
		return err
	}

	return copyDir(srcDir, dstDir)
}

// Step8CopyTLS copies TLS files from _output to ./
type Step8CopyTLS struct {
	*BaseStep
}

func NewStep8(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step8CopyTLS, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step8CopyTLS{BaseStep: base}, nil
}

func (s *Step8CopyTLS) Name() string {
	return "Copy TLS files"
}

func (s *Step8CopyTLS) Execute() error {
	srcDir := filepath.Join(s.cfg.OutputDir, "tls")
	dstDir := "tls"

	if err := util.EnsureDir(dstDir); err != nil {
		return err
	}

	return copyDir(srcDir, dstDir)
}

// Step9DeployCluster runs openshift-install create cluster
type Step9DeployCluster struct {
	*BaseStep
}

func NewStep9(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step9DeployCluster, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step9DeployCluster{BaseStep: base}, nil
}

func (s *Step9DeployCluster) Name() string {
	return "Deploy cluster"
}

func (s *Step9DeployCluster) Execute() error {
	versionDir := filepath.Join("artifacts", s.versionArch)
	installBin := util.GetBinaryPath(s.versionArch, "openshift-install")
	args := []string{"create", "cluster", "--dir", versionDir, "--log-level=debug"}

	return util.RunCommand(s.executor, installBin, args...)
}

// Step10Verify performs post-install verification
type Step10Verify struct {
	*BaseStep
}

func NewStep10(cfg *config.Config, log *logger.Logger, executor util.CommandExecutor) (*Step10Verify, error) {
	base, err := newBaseStep(cfg, log, executor)
	if err != nil {
		return nil, err
	}
	return &Step10Verify{BaseStep: base}, nil
}

func (s *Step10Verify) Name() string {
	return "Verify installation"
}

func (s *Step10Verify) Execute() error {
	// Check 1: Root credentials should not exist
	_, err := s.executor.Execute("oc", "get", "secrets", "-n", "kube-system", "aws-creds")
	if err == nil {
		s.log.Error("WARNING: Root credentials secret exists (expected it to not exist)")
	} else {
		s.log.Info("✓ Root credentials secret does not exist (as expected)")
	}

	// Check 2: Components should use IAM roles
	output, err := s.executor.Execute("oc", "get", "secrets", "-n", "openshift-image-registry",
		"installer-cloud-credentials", "-o", "json")
	if err != nil {
		return fmt.Errorf("failed to check IAM role usage: %w", err)
	}

	if len(output) > 0 && (contains(output, "role_arn") || contains(output, "web_identity_token_file")) {
		s.log.Info("✓ Components are using IAM roles")
	} else {
		s.log.Error("WARNING: Components may not be using IAM roles correctly")
	}

	return nil
}

// Helper function to copy directories
func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && haystack != "" && needle != "" &&
		findSubstring(haystack, needle)
}

func findSubstring(haystack, needle string) bool {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
