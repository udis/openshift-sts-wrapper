package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gitlab.cee.redhat.com/clobrano/ccoctl-sso/pkg/logger"
	"gitlab.cee.redhat.com/clobrano/ccoctl-sso/pkg/util"
	"github.com/spf13/cobra"
)

var (
	cleanupClusterName string
	cleanupAwsRegion   string
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up AWS resources after a failed installation",
	Long:  `Removes AWS resources (S3 bucket, IAM roles) created during installation`,
	Run:   runCleanup,
}

func init() {
	rootCmd.AddCommand(cleanupCmd)

	cleanupCmd.Flags().StringVar(&cleanupClusterName, "cluster-name", "", "Cluster/infrastructure name")
	cleanupCmd.Flags().StringVar(&cleanupAwsRegion, "region", "", "AWS region")
	cleanupCmd.MarkFlagRequired("cluster-name")
	cleanupCmd.MarkFlagRequired("region")
}

func runCleanup(cmd *cobra.Command, args []string) {
	log := logger.New(logger.Level(getLogLevel()), nil)

	// Confirm with user
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("This will delete AWS resources for cluster '%s' in region '%s'.\n", cleanupClusterName, cleanupAwsRegion)
	fmt.Print("Continue? (y/n): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" && response != "yes" {
		log.Info("Cleanup cancelled.")
		return
	}

	log.StartStep("Cleaning up AWS resources")

	executor := &util.RealExecutor{}

	// Find ccoctl binary - check common locations
	ccoctlPath := "ccoctl"
	if util.FileExists("artifacts/bin/ccoctl") {
		ccoctlPath = "artifacts/bin/ccoctl"
	}

	args_cleanup := []string{
		"aws", "delete",
		"--name", cleanupClusterName,
		"--region", cleanupAwsRegion,
	}

	if err := util.RunCommand(executor, ccoctlPath, args_cleanup...); err != nil {
		log.FailStep("Cleanup")
		log.Error(fmt.Sprintf("Failed to clean up: %v", err))
		log.Info("You may need to manually delete AWS resources.")
		os.Exit(1)
	}

	log.CompleteStep("Cleanup")
	log.Info("AWS resources have been deleted.")
}
