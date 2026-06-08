package kube

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tiamxu/kit/cli"
	"github.com/tiamxu/kit/log"
)

// Tool Kubernetes 工具
type Tool struct{}

// Name 工具名称
func (t *Tool) Name() string { return "kubectl" }

// Description 工具描述
func (t *Tool) Description() string { return "Manage kubernetes resources" }

// Flags 工具全局标志
func (t *Tool) Flags() []cli.Flag {
	// 不在根命令添加 Flags，只在子命令中添加，避免解析顺序问题
	return []cli.Flag{}
}

// Commands 工具命令
func (t *Tool) Commands() []*cli.Command {
	return []*cli.Command{
		cli.NewCommand("get").
			SetDescription("Get k8s resource deployment").
			AddFlags(cli.StringFlag("namespace", "n", "default", "Set kubernetes namespace")).
			AddFlags(cli.StringFlag("name", "d", "", "Set deployment name")).
			SetRun(func(ctx *cli.Context) error {
				return t.RunGetDeployment(ctx)
			}),
		cli.NewCommand("restart").
			SetDescription("Restart k8s resource deployment").
			AddFlags(cli.StringFlag("namespace", "n", "default", "Set kubernetes namespace")).
			AddFlags(cli.StringFlag("name", "d", "", "Set deployment name")).
			SetRun(func(ctx *cli.Context) error {
				return t.RunRestart(ctx)
			}),
		cli.NewCommand("create").
			SetDescription("Create resource deployment").
			AddFlags(cli.StringFlag("namespace", "n", "default", "Set kubernetes namespace")).
			AddFlags(cli.StringFlag("name", "d", "", "Set deployment name")).
			AddFlags(cli.StringFlag("image", "i", "", "Set deployment image")).
			AddFlags(cli.IntFlag("replicas", "r", 1, "Set number of replicas")).
			SetRun(func(ctx *cli.Context) error {
				return t.CreateDeployment(ctx)
			}),
	}
}

// RunGetDeployment 执行获取部署命令
func (t *Tool) RunGetDeployment(ctx *cli.Context) error {
	namespace := ctx.String("namespace")
	name := ctx.String("name")

	var kubectlCmd *exec.Cmd
	if name != "" {
		kubectlCmd = exec.Command("kubectl", "get", "deployment", name, "-n", namespace, "-o", "yaml")
	} else {
		kubectlCmd = exec.Command("kubectl", "get", "deployments", "-n", namespace)
	}

	log.Infof("Getting deployment in namespace %s", namespace)
	if err := t.executeCommand(kubectlCmd); err != nil {
		return fmt.Errorf("kubectl get deployment failed: %v", err)
	}

	return nil
}

// RunRestart 执行重启部署命令
func (t *Tool) RunRestart(ctx *cli.Context) error {
	namespace := ctx.String("namespace")
	name := ctx.String("name")

	if name == "" {
		return fmt.Errorf("deployment name is required")
	}

	kubectlCmd := exec.Command("kubectl", "rollout", "restart", "deployment", name, "-n", namespace)

	log.Infof("Restarting deployment %s in namespace %s", name, namespace)
	if err := t.executeCommand(kubectlCmd); err != nil {
		return fmt.Errorf("kubectl rollout restart failed: %v", err)
	}

	// 等待重启完成
	kubectlStatusCmd := exec.Command("kubectl", "rollout", "status", "deployment", name, "-n", namespace)
	log.Infof("Waiting for deployment %s to restart", name)
	if err := t.executeCommand(kubectlStatusCmd); err != nil {
		return fmt.Errorf("kubectl rollout status failed: %v", err)
	}

	log.Infof("Deployment %s restarted successfully", name)
	return nil
}

// CreateDeployment 执行创建部署命令
func (t *Tool) CreateDeployment(ctx *cli.Context) error {
	namespace := ctx.String("namespace")
	name := ctx.String("name")
	image := ctx.String("image")
	replicas := ctx.Int("replicas")

	if name == "" {
		return fmt.Errorf("deployment name is required")
	}
	if image == "" {
		return fmt.Errorf("container image is required")
	}

	kubectlCmd := exec.Command("kubectl", "create", "deployment", name, "--image", image, "--replicas", fmt.Sprintf("%d", replicas), "-n", namespace)

	log.Infof("Creating deployment %s with image %s and %d replicas in namespace %s", name, image, replicas, namespace)
	if err := t.executeCommand(kubectlCmd); err != nil {
		return fmt.Errorf("kubectl create deployment failed: %v", err)
	}

	log.Infof("Deployment %s created successfully", name)
	return nil
}

// executeCommand 执行命令
func (t *Tool) executeCommand(cmd *exec.Cmd) error {
	// 设置命令的标准输出和标准错误
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	return cmd.Run()
}
