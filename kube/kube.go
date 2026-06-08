package kube

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/tiamxu/kit/log"
)

// Tool Kubernetes 工具
type Tool struct{}

// NewTool 创建 Kubernetes 工具实例
func NewTool() *Tool { return &Tool{} }

// AddCommands 注册 Kubernetes 命令到根命令
func (t *Tool) AddCommands(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "kube",
		Short: "Manage kubernetes resources",
	}
	cmd.AddCommand(t.getCmd(), t.restartCmd(), t.createCmd())
	root.AddCommand(cmd)
}

func (t *Tool) getCmd() *cobra.Command {
	var namespace, name string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get k8s resource deployment",
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.RunGetDeployment(namespace, name)
		},
	}
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Set kubernetes namespace")
	cmd.Flags().StringVarP(&name, "name", "d", "", "Set deployment name")
	return cmd
}

func (t *Tool) restartCmd() *cobra.Command {
	var namespace, name string
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart k8s resource deployment",
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.RunRestart(namespace, name)
		},
	}
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Set kubernetes namespace")
	cmd.Flags().StringVarP(&name, "name", "d", "", "Set deployment name")
	return cmd
}

func (t *Tool) createCmd() *cobra.Command {
	var namespace, name, image string
	var replicas int
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create resource deployment",
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.CreateDeployment(namespace, name, image, replicas)
		},
	}
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Set kubernetes namespace")
	cmd.Flags().StringVarP(&name, "name", "d", "", "Set deployment name")
	cmd.Flags().StringVarP(&image, "image", "i", "", "Set deployment image")
	cmd.Flags().IntVarP(&replicas, "replicas", "r", 1, "Set number of replicas")
	return cmd
}

// RunGetDeployment 执行获取部署命令
func (t *Tool) RunGetDeployment(namespace, name string) error {
	var kubectlCmd *exec.Cmd
	if name != "" {
		kubectlCmd = exec.Command("kubectl", "get", "deployment", name, "-n", namespace, "-o", "yaml")
	} else {
		kubectlCmd = exec.Command("kubectl", "get", "deployments", "-n", namespace)
	}

	log.Infof("Getting deployment in namespace %s", namespace)
	if err := t.executeCommand(kubectlCmd); err != nil {
		return fmt.Errorf("kubectl get deployment failed: %w", err)
	}

	return nil
}

// RunRestart 执行重启部署命令
func (t *Tool) RunRestart(namespace, name string) error {
	if name == "" {
		return fmt.Errorf("deployment name is required")
	}

	kubectlCmd := exec.Command("kubectl", "rollout", "restart", "deployment", name, "-n", namespace)

	log.Infof("Restarting deployment %s in namespace %s", name, namespace)
	if err := t.executeCommand(kubectlCmd); err != nil {
		return fmt.Errorf("kubectl rollout restart failed: %w", err)
	}

	kubectlStatusCmd := exec.Command("kubectl", "rollout", "status", "deployment", name, "-n", namespace)
	log.Infof("Waiting for deployment %s to restart", name)
	if err := t.executeCommand(kubectlStatusCmd); err != nil {
		return fmt.Errorf("kubectl rollout status failed: %w", err)
	}

	log.Infof("Deployment %s restarted successfully", name)
	return nil
}

// CreateDeployment 执行创建部署命令
func (t *Tool) CreateDeployment(namespace, name, image string, replicas int) error {
	if name == "" {
		return fmt.Errorf("deployment name is required")
	}
	if image == "" {
		return fmt.Errorf("container image is required")
	}

	kubectlCmd := exec.Command("kubectl", "create", "deployment", name, "--image", image, "--replicas", fmt.Sprintf("%d", replicas), "-n", namespace)

	log.Infof("Creating deployment %s with image %s and %d replicas in namespace %s", name, image, replicas, namespace)
	if err := t.executeCommand(kubectlCmd); err != nil {
		return fmt.Errorf("kubectl create deployment failed: %w", err)
	}

	log.Infof("Deployment %s created successfully", name)
	return nil
}

func (t *Tool) executeCommand(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
