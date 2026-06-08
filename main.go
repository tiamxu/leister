package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tiamxu/kit/log"
	"github.com/tiamxu/leister/client"
	"github.com/tiamxu/leister/config"
	"github.com/tiamxu/leister/docker"
	"github.com/tiamxu/leister/gitlab"
	"github.com/tiamxu/leister/jenkins"
	"github.com/tiamxu/leister/kube"
)

func main() {
	if err := log.InitLogger(&log.Config{
		Level:  "info",
		Type:   "stdout",
		Format: "text",
	}); err != nil {
		panic(err)
	}
	defer log.Sync()

	cfg := config.Load()
	apiClient := client.NewClient(cfg)

	rootCmd := &cobra.Command{
		Use:     "gigctl",
		Short:   "DevOps tools for building, CI/CD and deployment",
		Long:    "gigctl is a CLI toolset for Docker, Kubernetes, Jenkins and GitLab operations.",
		Version: "0.0.1",
	}

	docker.NewTool().AddCommands(rootCmd)
	kube.NewTool().AddCommands(rootCmd)
	jenkins.NewTool(apiClient).AddCommands(rootCmd)
	gitlab.NewTool(apiClient).AddCommands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
