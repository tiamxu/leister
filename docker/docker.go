package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tiamxu/kit/log"
)

const (
	RegistryDomain    = "harbor.xuliang.net.cn"
	RegistryNamespace = "xuliang"
	DefaultEnv        = "dev"
	// 登录账号从环境变量读取，缺失即报错
	// 使用：DOCKER_REGISTRY_USERNAME / DOCKER_REGISTRY_PASSWORD
	RegistryUsernameEnv = "DOCKER_REGISTRY_USERNAME"
	RegistryPasswordEnv = "DOCKER_REGISTRY_PASSWORD"
)

// Tool Docker 工具
type Tool struct{}

// NewTool 创建 Docker 工具实例
func NewTool() *Tool { return &Tool{} }

// AddCommands 注册 Docker 命令到根命令
func (t *Tool) AddCommands(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "docker",
		Short: "Manage docker commands",
	}
	cmd.AddCommand(t.buildCmd(), t.pushCmd())
	root.AddCommand(cmd)
}

func (t *Tool) buildCmd() *cobra.Command {
	var tag, env, lang, dockerfile string
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build an image from a Dockerfile",
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.RunBuild(tag, env, lang, dockerfile)
		},
	}
	cmd.Flags().StringVarP(&tag, "tag", "t", "", "Set docker image tag")
	cmd.Flags().StringVarP(&env, "env", "e", "dev", "Set docker image env")
	cmd.Flags().StringVarP(&lang, "lang", "l", "go", "Set type of code language")
	cmd.Flags().StringVarP(&dockerfile, "dockerfile", "f", "Dockerfile-dev", "Path to Dockerfile")
	return cmd
}

func (t *Tool) pushCmd() *cobra.Command {
	var tag, env string
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Upload an image to a registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.RunPush(tag, env)
		},
	}
	cmd.Flags().StringVarP(&tag, "tag", "t", "", "Set docker image tag")
	cmd.Flags().StringVarP(&env, "env", "e", "dev", "Set docker image env")
	return cmd
}

// RunBuild 执行构建命令
func (t *Tool) RunBuild(tag, env, lang, dockerfile string) error {
	if err := t.loginAction(); err != nil {
		return err
	}
	return t.buildAction(tag, lang, dockerfile)
}

// RunPush 执行推送命令
func (t *Tool) RunPush(tag, env string) error {
	if err := t.loginAction(); err != nil {
		return err
	}
	return t.pushAction(tag)
}

// loginAction 执行 Docker 登录操作
func (t *Tool) loginAction() error {
	username := os.Getenv(RegistryUsernameEnv)
	password := os.Getenv(RegistryPasswordEnv)
	if username == "" || password == "" {
		return fmt.Errorf("请设置环境变量 %s 和 %s", RegistryUsernameEnv, RegistryPasswordEnv)
	}

	registryURL := "https://" + RegistryDomain

	var loginCmd *exec.Cmd
	if t.needSudo() {
		loginCmd = exec.Command("sudo", "docker", "login",
			"--username="+username, "--password="+password,
			registryURL)
	} else {
		loginCmd = exec.Command("docker", "login",
			"--username="+username, "--password="+password,
			registryURL)
	}

	log.Infof("Logging in to Docker registry: %s", registryURL)
	if err := t.executeCommand(loginCmd); err != nil {
		return fmt.Errorf("docker login failed: %w", err)
	}

	log.Infof("Docker login success")
	return nil
}

// buildAction 执行构建操作
func (t *Tool) buildAction(tag, lang, dockerfile string) error {
	if tag == "" {
		tag = time.Now().Format("200601021504")
	}

	projectCtx, err := t.getProjectContext()
	if err != nil {
		return err
	}

	if err := t.buildCode(lang); err != nil {
		return err
	}

	imageName := fmt.Sprintf("%s/%s/%s:%s", RegistryDomain, RegistryNamespace, projectCtx, tag)
	dockerBuildCmd := t.getDockerBuildCommand(dockerfile, imageName)
	log.Infof("Building Docker image: %s", imageName)
	if err := t.executeCommand(dockerBuildCmd); err != nil {
		return fmt.Errorf("docker build failed: %w", err)
	}

	log.Infof("Docker build success: %s", imageName)
	return nil
}

// pushAction 执行推送操作
func (t *Tool) pushAction(tag string) error {
	if tag == "" {
		tag = time.Now().Format("200601021504")
	}

	projectCtx, err := t.getProjectContext()
	if err != nil {
		return err
	}

	imageName := fmt.Sprintf("%s/%s/%s:%s", RegistryDomain, RegistryNamespace, projectCtx, tag)
	dockerPushCmd := t.getDockerPushCommand(imageName)
	log.Infof("Pushing Docker image: %s", imageName)
	if err := t.executeCommand(dockerPushCmd); err != nil {
		return fmt.Errorf("docker push failed: %w", err)
	}

	log.Infof("Docker push success: %s", imageName)
	return nil
}

// getProjectContext 获取项目上下文
func (t *Tool) getProjectContext() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	return filepath.Base(pwd), nil
}

// buildCode 构建代码
func (t *Tool) buildCode(lang string) error {
	switch lang {
	case "go":
		goBuildCmd := exec.Command("go", "build", ".")
		log.Infof("Building Go code")
		if err := t.executeCommand(goBuildCmd); err != nil {
			return fmt.Errorf("go build failed: %w", err)
		}
	case "node":
		npmInstallCmd := exec.Command("npm", "install")
		log.Infof("Installing Node.js dependencies")
		if err := t.executeCommand(npmInstallCmd); err != nil {
			return fmt.Errorf("npm install failed: %w", err)
		}

		npmBuildCmd := exec.Command("npm", "run", "build")
		log.Infof("Building Node.js code")
		if err := t.executeCommand(npmBuildCmd); err != nil {
			return fmt.Errorf("npm run build failed: %w", err)
		}
	default:
		return fmt.Errorf("unsupported language: %s", lang)
	}

	return nil
}

func (t *Tool) getDockerBuildCommand(dockerfile, imageName string) *exec.Cmd {
	if t.needSudo() {
		return exec.Command("sudo", "docker", "build", "-f", dockerfile, "-t", imageName, ".")
	}
	return exec.Command("docker", "build", "-f", dockerfile, "-t", imageName, ".")
}

func (t *Tool) getDockerPushCommand(imageName string) *exec.Cmd {
	if t.needSudo() {
		return exec.Command("sudo", "docker", "push", imageName)
	}
	return exec.Command("docker", "push", imageName)
}

func (t *Tool) needSudo() bool {
	if runtime.GOOS == "windows" {
		return false
	}
	dockerGroupCmd := exec.Command("groups")
	output, err := dockerGroupCmd.CombinedOutput()
	if err != nil {
		return true
	}
	return !strings.Contains(string(output), "docker")
}

func (t *Tool) executeCommand(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
