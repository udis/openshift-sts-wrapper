package steps

import (
	"gitlab.cee.redhat.com/clobrano/ccoctl-sso/pkg/config"
	"gitlab.cee.redhat.com/clobrano/ccoctl-sso/pkg/util"
)

type Detector struct {
	cfg         *config.Config
	versionArch string
}

func NewDetector(cfg *config.Config) *Detector {
	versionArch, _ := util.ExtractVersionArch(cfg.ReleaseImage)
	return &Detector{
		cfg:         cfg,
		versionArch: versionArch,
	}
}

func (d *Detector) ShouldSkipStep(stepNum int) bool {
	// If StartFromStep is set, skip all steps before it
	if d.cfg.StartFromStep > 0 && stepNum < d.cfg.StartFromStep {
		return true
	}

	// Otherwise, check for evidence of completion
	switch stepNum {
	case 1:
		// Step 1: Extract credentials requests
		return util.DirExistsWithFiles(util.GetCredReqsPath(d.versionArch))
	case 2:
		// Step 2: Extract binaries
		return util.FileExists(util.GetBinaryPath(d.versionArch, "openshift-install")) &&
			util.FileExists(util.GetBinaryPath(d.versionArch, "ccoctl"))
	case 3:
		// Step 3: Create install-config.yaml
		return util.FileExists(util.GetInstallConfigPath(d.versionArch))
	case 4:
		// Step 4: Set credentialsMode
		return util.FileContains(util.GetInstallConfigPath(d.versionArch), "credentialsMode: Manual")
	case 5:
		// Step 5: Create manifests
		return util.DirExistsWithFiles("manifests")
	case 6:
		// Step 6: Create AWS resources
		return util.DirExistsWithFiles("_output/manifests") &&
			util.DirExistsWithFiles("_output/tls")
	case 7:
		// Step 7: Copy manifests
		return util.DirExistsWithFiles("manifests")
	case 8:
		// Step 8: Copy TLS
		return util.DirExistsWithFiles("tls")
	case 9:
		// Step 9: Deploy cluster
		return util.FileExists(".openshift_install.log")
	case 10:
		// Step 10: Verify installation
		// Verification should always run, don't skip it
		return false
	default:
		return false
	}
}
