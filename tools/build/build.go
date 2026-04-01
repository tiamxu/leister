package build

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tiamxu/kit/cli"
	"github.com/tiamxu/kit/log"
)

const (
	RegistryDomain    = "registry.cn-hangzhou.aliyuncs.com"
	RegistryNamespace = "unipal"
	Env               = "dev"
)

type Tool struct{}

func (t *Tool) Name() string        { return "deploy" }
func (t *Tool) Description() string { return "Manager deploy server" }

func (t *Tool) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag("tag", "t", "latest", "Set docker image tag"),
		cli.StringFlag("env", "e", "dev", "Set docker image env"),
		cli.StringFlag("lang", "l", "go", "Set type of code langue"),
	}
}

func (t *Tool) Commands() []*cli.Command {
	return []*cli.Command{
		cli.NewCommand("build").
			SetDescription("Build code and docker image").
			SetRun(func(ctx *cli.Context) error {
				return RunBuild(ctx)
			}),
		cli.NewCommand("push").
			SetDescription("Docker push image registry").
			SetRun(func(ctx *cli.Context) error {
				return RunPush(ctx)
			}),
	}
}

type Context struct {
	Name, AbsPath string
}

func Initial(c *cli.Context) (*Context, error) {
	absPath, err := filepath.Abs("./")
	if err != nil {
		return nil, errors.New("get project absolute path failure, reason: " + err.Error())
	}
	absPath = strings.Replace(absPath, `\`, `/`, -1)
	absPath = strings.TrimRight(absPath, "/")
	var (
		ctx = &Context{
			Name:    filepath.Base(absPath),
			AbsPath: absPath,
		}
	)
	return ctx, nil
}

func needSudo() bool {
	sysType := runtime.GOOS
	if sysType == "linux" {
		if os.Getuid() != 0 {
			return true
		}
	}
	return false
}

func RunBuild(ctx *cli.Context) error {
	if err := loginAction(ctx); err != nil {
		return err
	}
	return buildAction(ctx)
}

func RunPush(ctx *cli.Context) error {
	if err := loginAction(ctx); err != nil {
		return err
	}
	return pushAction(ctx)
}

func buildAction(ctx *cli.Context) error {
	// 检查 Dockerfile
	dockerfile := "./Dockerfile-dev"
	if _, err := os.Stat(dockerfile); os.IsNotExist(err) {
		log.Errorf("Dockerfile not found: %s", dockerfile)
		return errors.New("docker build failed, not found Dockerfile")
	}
	log.Infof("Found Dockerfile: %s", dockerfile)

	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)

	// 代码构建
	lang := ctx.String("lang")
	env := ctx.String("env")
	log.Infof("Building with language: %s, environment: %s", lang, env)

	var codeBuilder *exec.Cmd
	switch lang {
	case "node":
		// 先安装依赖
		log.Infof("Installing Node.js dependencies...")
		npmInstall := exec.Command("npm", "install", "--registry=https://registry.npm.taobao.org")
		npmInstall.Stderr = &stderr
		npmInstall.Stdout = &stdout
		if err := npmInstall.Run(); err != nil {
			log.Errorf("Failed to install dependencies: %s", stderr.String())
			return fmt.Errorf("npm install failed: %w", err)
		}
		log.Infof("Dependencies installed successfully")
		
		// 构建代码
		buildCmd := "npm run build" + ":" + env
		codeBuilder = exec.Command("bash", "-c", buildCmd)
	case "go":
		// 确保 bin 目录存在
		if err := os.MkdirAll("bin", 0755); err != nil {
			log.Errorf("Failed to create bin directory: %v", err)
			return fmt.Errorf("create bin directory failed: %w", err)
		}
		codeBuilder = exec.Command("go", "build", "-o", "bin/main")
		codeBuilder.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")
	default:
		return fmt.Errorf("unsupported language: %s", lang)
	}

	// 执行代码构建
	codeBuilder.Stderr = &stderr
	codeBuilder.Stdout = &stdout
	log.Infof("Executing build command: %s", codeBuilder.String())
	if err := codeBuilder.Run(); err != nil {
		log.Errorf("Build failed: %s", stderr.String())
		return fmt.Errorf("code build failed: %w", err)
	}
	log.Infof("Build output: %s", stdout.String())
	log.Infof("%s build completed successfully", lang)

	// Docker 镜像构建
	ctxBuild, err := Initial(ctx)
	if err != nil {
		log.Errorf("Failed to initialize build context: %v", err)
		return fmt.Errorf("init build context failed: %w", err)
	}

	version := ctx.String("tag")
	if version == "" {
		version = "latest"
		log.Warnf("Tag not specified, using default: latest")
	}

	registryPath := fmt.Sprintf("%s/%s/%s:%s", RegistryDomain, RegistryNamespace, strings.ToLower(env+"_"+ctxBuild.Name), version)
	log.Infof("Building Docker image: %s", registryPath)

	var dockerBuilder *exec.Cmd
	args := []string{
		"build",
		"-f", dockerfile,
		"-t", registryPath,
		"--platform", "linux/amd64", // 确保构建 Linux 镜像
		".",
	}

	if needSudo() {
		dockerBuilder = exec.Command("sudo", append([]string{"docker"}, args...)...)
	} else {
		dockerBuilder = exec.Command("docker", args...)
	}

	dockerBuilder.Stderr = &stderr
	dockerBuilder.Stdout = &stdout
	log.Infof("Executing Docker build: %s", dockerBuilder.String())
	if err := dockerBuilder.Run(); err != nil {
		log.Errorf("Docker build failed: %s", stderr.String())
		return fmt.Errorf("docker build failed: %w", err)
	}
	log.Infof("Docker build output: %s", stdout.String())
	log.Infof("Docker image built successfully: %s", registryPath)

	return nil
}

func pushAction(ctx *cli.Context) error {
	var (
		pushCmd *exec.Cmd
		stderr  bytes.Buffer
		stdout  bytes.Buffer
	)

	ctxBuild, err := Initial(ctx)
	if err != nil {
		log.Errorf("Failed to initialize build context: %v", err)
		return fmt.Errorf("init build context failed: %w", err)
	}

	version := ctx.String("tag")
	if version == "" {
		version = "latest"
		log.Warnf("Tag not specified, using default: latest")
	}

	env := ctx.String("env")
	if env == "" {
		env = Env
		log.Warnf("Environment not specified, using default: %s", Env)
	}

	registryPath := fmt.Sprintf("%s/%s/%s:%s", RegistryDomain, RegistryNamespace, strings.ToLower(env+"_"+ctxBuild.Name), version)
	log.Infof("Pushing Docker image: %s", registryPath)

	args := []string{"push", registryPath}

	if needSudo() {
		pushCmd = exec.Command("sudo", append([]string{"docker"}, args...)...)
	} else {
		pushCmd = exec.Command("docker", args...)
	}

	pushCmd.Stderr = &stderr
	pushCmd.Stdout = &stdout
	log.Infof("Executing Docker push: %s", pushCmd.String())

	if err := pushCmd.Run(); err != nil {
		log.Errorf("Docker push failed: %s", stderr.String())
		return fmt.Errorf("docker push failed: %w", err)
	}

	log.Infof("Docker push output: %s", stdout.String())
	log.Infof("Docker image pushed successfully: %s", registryPath)

	return nil
}

func loginAction(ctx *cli.Context) error {
	var (
		loginCmd *exec.Cmd
		stderr   bytes.Buffer
		stdout   bytes.Buffer
	)

	registry := "https://registry.cn-hangzhou.aliyuncs.com"
	username := "root"
	password := "123456"

	log.Infof("Logging into Docker registry: %s", registry)

	args := []string{
		"login",
		"--username", username,
		"--password", password,
		registry,
	}

	if needSudo() {
		loginCmd = exec.Command("sudo", append([]string{"docker"}, args...)...)
	} else {
		loginCmd = exec.Command("docker", args...)
	}

	loginCmd.Stderr = &stderr
	loginCmd.Stdout = &stdout
	log.Infof("Executing Docker login: %s", loginCmd.String())

	if err := loginCmd.Run(); err != nil {
		log.Errorf("Docker login failed: %s", stderr.String())
		return fmt.Errorf("docker login failed: %w", err)
	}

	log.Infof("Docker login output: %s", stdout.String())
	log.Infof("Docker login successful")

	return nil
}
