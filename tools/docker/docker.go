package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/tiamxu/kit/cli"
	"github.com/tiamxu/kit/log"
)

const (
	RegistryDomain    = "harbor.yeemiao.net.cn"
	RegistryNamespace = "xuliang"
	DefaultEnv        = "dev"
	// 登录信息硬编码
	RegistryUsername = "xuliang"
	RegistryPassword = "nD!dfjk1s613dv"
)

// Tool Docker 工具
type Tool struct{}

// Name 工具名称
func (t *Tool) Name() string { return "docker" }

// Description 工具描述
func (t *Tool) Description() string { return "Manage docker commands" }

// Flags 工具全局标志
func (t *Tool) Flags() []cli.Flag {
	// 不在根命令添加 Flags，只在子命令中添加，避免解析顺序问题
	return []cli.Flag{}
}

// Commands 工具命令
func (t *Tool) Commands() []*cli.Command {
	return []*cli.Command{
		cli.NewCommand("build").
			SetDescription("Build an image from a Dockerfile").
			AddFlags(cli.StringFlag("tag", "t", "", "Set docker image tag")).
			AddFlags(cli.StringFlag("env", "e", "dev", "Set docker image env")).
			AddFlags(cli.StringFlag("lang", "l", "go", "Set type of code language")).
			AddFlags(cli.StringFlag("dockerfile", "f", "Dockerfile-dev", "Path to Dockerfile")).
			SetRun(func(ctx *cli.Context) error {
				return t.RunBuild(ctx)
			}),
		cli.NewCommand("push").
			SetDescription("Upload an image to a registry").
			AddFlags(cli.StringFlag("tag", "t", "", "Set docker image tag")).
			AddFlags(cli.StringFlag("env", "e", "dev", "Set docker image env")).
			SetRun(func(ctx *cli.Context) error {
				return t.RunPush(ctx)
			}),
	}
}

// RunBuild 执行构建命令
func (t *Tool) RunBuild(ctx *cli.Context) error {
	// 登录到 Docker registry
	if err := t.loginAction(); err != nil {
		return err
	}

	// 执行构建操作
	return t.buildAction(ctx)
}

// RunPush 执行推送命令
func (t *Tool) RunPush(ctx *cli.Context) error {
	// 登录到 Docker registry
	if err := t.loginAction(); err != nil {
		return err
	}

	// 执行推送操作
	return t.pushAction(ctx)
}

// loginAction 执行 Docker 登录操作
func (t *Tool) loginAction() error {
	registryURL := "https://" + RegistryDomain

	var loginCmd *exec.Cmd
	if t.needSudo() {
		loginCmd = exec.Command("sudo", "docker", "login",
			"--username="+RegistryUsername, "--password="+RegistryPassword,
			registryURL)
	} else {
		loginCmd = exec.Command("docker", "login",
			"--username="+RegistryUsername, "--password="+RegistryPassword,
			registryURL)
	}

	log.Infof("Logging in to Docker registry: %s", registryURL)
	if err := t.executeCommand(loginCmd); err != nil {
		return fmt.Errorf("docker login failed: %v", err)
	}

	log.Infof("Docker login success")
	return nil
}

// buildAction 执行构建操作
func (t *Tool) buildAction(ctx *cli.Context) error {
	// 获取参数
	tag := ctx.String("tag")
	if tag == "" {
		tag = time.Now().Format("200601021504")
	}
	// env := ctx.String("env") // 暂时未使用
	lang := ctx.String("lang")
	dockerfile := ctx.String("dockerfile")

	// 获取项目上下文
	projectCtx, err := t.getProjectContext()
	if err != nil {
		return err
	}

	// 构建代码
	if err := t.buildCode(lang); err != nil {
		return err
	}

	// 构建 Docker 镜像
	imageName := fmt.Sprintf("%s/%s/%s:%s", RegistryDomain, RegistryNamespace, projectCtx, tag)
	dockerBuildCmd := t.getDockerBuildCommand(dockerfile, imageName)
	log.Infof("Building Docker image: %s", imageName)
	if err := t.executeCommand(dockerBuildCmd); err != nil {
		return fmt.Errorf("docker build failed: %v", err)
	}

	log.Infof("Docker build success: %s", imageName)
	return nil
}

// pushAction 执行推送操作
func (t *Tool) pushAction(ctx *cli.Context) error {
	// 获取参数
	tag := ctx.String("tag")
	if tag == "" {
		tag = time.Now().Format("200601021504")
	}

	// 获取项目上下文
	projectCtx, err := t.getProjectContext()
	if err != nil {
		return err
	}

	// 推送 Docker 镜像
	imageName := fmt.Sprintf("%s/%s/%s:%s", RegistryDomain, RegistryNamespace, projectCtx, tag)
	dockerPushCmd := t.getDockerPushCommand(imageName)
	log.Infof("Pushing Docker image: %s", imageName)
	if err := t.executeCommand(dockerPushCmd); err != nil {
		return fmt.Errorf("docker push failed: %v", err)
	}

	log.Infof("Docker push success: %s", imageName)
	return nil
}

// getProjectContext 获取项目上下文
func (t *Tool) getProjectContext() (string, error) {
	// 获取当前工作目录
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// 获取当前目录名称作为项目上下文
	projectCtx := filepath.Base(pwd)
	return projectCtx, nil
}

// buildCode 构建代码
func (t *Tool) buildCode(lang string) error {
	switch lang {
	case "go":
		// 构建 Go 代码
		goBuildCmd := exec.Command("go", "build", ".")
		log.Infof("Building Go code")
		if err := t.executeCommand(goBuildCmd); err != nil {
			return fmt.Errorf("go build failed: %v", err)
		}
	case "node":
		// 构建 Node.js 代码
		npmInstallCmd := exec.Command("npm", "install")
		log.Infof("Installing Node.js dependencies")
		if err := t.executeCommand(npmInstallCmd); err != nil {
			return fmt.Errorf("npm install failed: %v", err)
		}

		npmBuildCmd := exec.Command("npm", "run", "build")
		log.Infof("Building Node.js code")
		if err := t.executeCommand(npmBuildCmd); err != nil {
			return fmt.Errorf("npm run build failed: %v", err)
		}
	default:
		return fmt.Errorf("unsupported language: %s", lang)
	}

	return nil
}

// getDockerBuildCommand 获取 Docker 构建命令
func (t *Tool) getDockerBuildCommand(dockerfile, imageName string) *exec.Cmd {
	if t.needSudo() {
		return exec.Command("sudo", "docker", "build", "-f", dockerfile, "-t", imageName, ".")
	}
	return exec.Command("docker", "build", "-f", dockerfile, "-t", imageName, ".")
}

// getDockerPushCommand 获取 Docker 推送命令
func (t *Tool) getDockerPushCommand(imageName string) *exec.Cmd {
	if t.needSudo() {
		return exec.Command("sudo", "docker", "push", imageName)
	}
	return exec.Command("docker", "push", imageName)
}

// needSudo 检查是否需要 sudo
func (t *Tool) needSudo() bool {
	// Windows 不需要 sudo
	if runtime.GOOS == "windows" {
		return false
	}

	// 检查当前用户是否在 docker 组
	dockerGroupCmd := exec.Command("groups")
	output, err := dockerGroupCmd.CombinedOutput()
	if err != nil {
		return true
	}

	return !strings.Contains(string(output), "docker")
}

// executeCommand 执行命令
func (t *Tool) executeCommand(cmd *exec.Cmd) error {
	// 设置命令的标准输出和标准错误
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	return cmd.Run()
}
