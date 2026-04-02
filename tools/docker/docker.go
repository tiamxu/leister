package docker

import (
	"bytes"
	"errors"
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
	RegistryNamespace = "yeemiao"
	DefaultEnv        = "dev"
	// 登录信息硬编码
	RegistryUsername = "xuliang"
	RegistryPassword = "nD!dfjk1s613dv"
)

type Tool struct{}

func (t *Tool) Name() string        { return "docker" }
func (t *Tool) Description() string { return "Manager docker cmd" }

func (t *Tool) Flags() []cli.Flag {
	// 不在根命令添加 Flags，只在子命令中添加，避免解析顺序问题
	return []cli.Flag{}
}

func (t *Tool) Commands() []*cli.Command {
	return []*cli.Command{
		cli.NewCommand("build").
			SetDescription("Build an image from a Dockerfile").
			AddFlags(
				cli.StringFlag("tag", "t", "latest", "Set docker image tag"),
				cli.StringFlag("env", "e", "dev", "Set docker image env"),
				cli.StringFlag("lang", "l", "go", "Set type of code langue"),
				cli.StringFlag("dockerfile", "f", "Dockerfile", "Path to Dockerfile"),
			).
			SetRun(func(ctx *cli.Context) error {
				return t.RunBuild(ctx)
			}),
		cli.NewCommand("push").
			SetDescription("Upload an image to a registry").
			AddFlags(
				cli.StringFlag("tag", "t", "latest", "Set docker image tag"),
				cli.StringFlag("env", "e", "dev", "Set docker image env"),
			).
			SetRun(func(ctx *cli.Context) error {
				return t.RunPush(ctx)
			}),
	}
}

type ProjectContext struct {
	Name, AbsPath string
}

func (t *Tool) getProjectContext() (*ProjectContext, error) {
	absPath, err := filepath.Abs("./")
	if err != nil {
		return nil, errors.New("get project absolute path failure, reason: " + err.Error())
	}
	absPath = strings.Replace(absPath, `\`, `/`, -1)
	absPath = strings.TrimRight(absPath, "/")
	return &ProjectContext{
		Name:    filepath.Base(absPath),
		AbsPath: absPath,
	}, nil
}

func (t *Tool) needSudo() bool {
	sysType := runtime.GOOS
	if sysType == "linux" {
		if os.Getuid() != 0 {
			return true
		}
	}
	return false
}

func (t *Tool) executeCommand(cmd *exec.Cmd) error {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	log.Infof("Executing: %s", cmd.String())
	if err := cmd.Run(); err != nil {
		log.Errorf("Command failed: %v", err)
		log.Errorf("Stderr: %s", stderr.String())
		return fmt.Errorf("command failed: %v, stderr: %s", err, stderr.String())
	}

	log.Infof("Command output: %s", stdout.String())
	return nil
}

func (t *Tool) RunBuild(ctx *cli.Context) error {
	if err := t.loginAction(); err != nil {
		return err
	}
	return t.buildAction(ctx)
}

func (t *Tool) RunPush(ctx *cli.Context) error {
	if err := t.loginAction(); err != nil {
		return err
	}
	return t.pushAction(ctx)
}

func (t *Tool) buildAction(ctx *cli.Context) error {
	dockerfile := ctx.String("dockerfile")
	// Check if Dockerfile exists
	if _, err := os.Stat(dockerfile); os.IsNotExist(err) {
		return errors.New("docker build failed, not found Dockerfile: " + dockerfile)
	}

	// Code language build
	lang := ctx.String("lang")
	env := ctx.String("env")
	log.Infof("Building with lang: %s, env: %s, dockerfile: %s", lang, env, dockerfile)

	if lang == "node" {
		// Install dependencies
		npmInstallCmd := exec.Command("npm", "install", "--registry=https://registry.npm.taobao.org")
		if err := t.executeCommand(npmInstallCmd); err != nil {
			return err
		}

		// Build project
		cmd := "npm run build:" + env
		npmBuildCmd := exec.Command("bash", "-c", cmd)
		if err := t.executeCommand(npmBuildCmd); err != nil {
			return err
		}
	} else {
		// Build Go project
		goBuildCmd := exec.Command("go", "build", "-o", "bin/main")
		goBuildCmd.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")
		if err := t.executeCommand(goBuildCmd); err != nil {
			return err
		}
	}

	log.Infof("%s build complete", lang)

	// Docker image build
	ctxBuild, err := t.getProjectContext()
	if err != nil {
		return err
	}

	version := ctx.String("tag")
	if version == "" {
		// 使用 202502031112 格式的时间戳
		version = time.Now().Format("200601021504")
	}

	registryPath := fmt.Sprintf("%s/%s/%s:%s", RegistryDomain, RegistryNamespace, strings.ToLower(env+"_"+ctxBuild.Name), version)

	var dockerBuildCmd *exec.Cmd
	if t.needSudo() {
		dockerBuildCmd = exec.Command("sudo", "docker", "build",
			"-f", dockerfile,
			"-t", registryPath, ".")
	} else {
		dockerBuildCmd = exec.Command("docker", "build",
			"-f", dockerfile,
			"-t", registryPath, ".")
	}

	if err := t.executeCommand(dockerBuildCmd); err != nil {
		return err
	}

	log.Infof("Docker build complete: %s", registryPath)
	return nil
}

func (t *Tool) pushAction(ctx *cli.Context) error {
	ctxBuild, err := t.getProjectContext()
	if err != nil {
		return err
	}

	version := ctx.String("tag")
	if version == "" {
		// 使用 202502031112 格式的时间戳
		version = time.Now().Format("200601021504")
	}

	env := ctx.String("env")
	if env == "" {
		env = DefaultEnv
	}

	registryPath := fmt.Sprintf("%s/%s/%s:%s", RegistryDomain, RegistryNamespace, strings.ToLower(env+"_"+ctxBuild.Name), version)

	var pushCmd *exec.Cmd
	if t.needSudo() {
		pushCmd = exec.Command("sudo", "docker", "push", registryPath)
	} else {
		pushCmd = exec.Command("docker", "push", registryPath)
	}

	log.Infof("Pushing Docker image: %s", registryPath)
	if err := t.executeCommand(pushCmd); err != nil {
		return err
	}

	log.Infof("Docker push complete: %s", registryPath)
	return nil
}

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
