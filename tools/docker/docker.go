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

	"github.com/tiamxu/kit/cli"
)

const (
	RegistryDomain    = "registry.cn-hangzhou.aliyuncs.com"
	RegistryNamespace = "unipal"
	Env               = "dev"
)

type Tool struct{}

func (t *Tool) Name() string        { return "docker" }
func (t *Tool) Description() string { return "Manager docker cmd" }

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
	//dockerfile
	if _, err := os.Stat("./Dockerfile-dev"); os.IsNotExist(err) {
		return errors.New("docker build failed,not found Dockerfile")
	}
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	//code langue build
	var codeBuilder *exec.Cmd
	lang := ctx.String("lang")
	env := ctx.String("env")
	fmt.Printf("lang:%s,env:%s\n", lang, env)
	if lang == "node" {
		cmd := "npm run build" + ":" + env
		exec.Command("npm", "install", "--registry=https://registry.npm.taobao.org")
		codeBuilder = exec.Command("bash", "-c", cmd)
	} else {
		codeBuilder = exec.Command("go", "build", "-o", "bin/main")
		codeBuilder.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")
	}
	codeBuilder.Stderr = &stderr
	codeBuilder.Stdout = &stdout
	fmt.Println(codeBuilder.String())
	if err := codeBuilder.Run(); err != nil {
		fmt.Println(stderr.String())
		return err
	}
	fmt.Println(stdout.String())
	fmt.Printf("%s build complate\n", ctx.String("lang"))
	//docker image build
	var dockerBuilder *exec.Cmd
	ctxBuild, _ := Initial(ctx)
	version := ctx.String("tag")
	if version == "" {
		version = "latest"
	}
	registryPath := fmt.Sprintf("%s/%s/%s:%s", RegistryDomain, RegistryNamespace, strings.ToLower(env+"_"+ctxBuild.Name), version)
	if needSudo() {
		dockerBuilder = exec.Command("sudo", "docker", "build",
			"-f", "Dockerfile-dev",
			"-t", registryPath, ".")
	} else {
		dockerBuilder = exec.Command("docker", "build",
			"-f", "Dockerfile-dev",
			"-t", registryPath, ".")
	}
	dockerBuilder.Stderr = &stderr
	dockerBuilder.Stdout = &stdout
	fmt.Println(dockerBuilder.String())
	if err := dockerBuilder.Run(); err != nil {
		fmt.Println(stderr.String())
		return err
	}
	fmt.Println("docker build complate")
	return nil
}

func pushAction(ctx *cli.Context) error {
	var (
		pushCmd *exec.Cmd
		stderr  bytes.Buffer
		stdout  bytes.Buffer
	)
	ctxBuild, _ := Initial(ctx)
	version := ctx.String("tag")
	if version == "" {
		version = "latest"
	}
	env := ctx.String("env")
	if env == "" {
		env = Env
	}
	registryPath := fmt.Sprintf("%s/%s/%s:%s", RegistryDomain, RegistryNamespace, strings.ToLower(env+"_"+ctxBuild.Name), version)
	if needSudo() {
		pushCmd = exec.Command("sudo", "docker", "push", registryPath)
	} else {
		pushCmd = exec.Command("docker", "push", registryPath)
	}
	pushCmd.Stderr = &stderr
	pushCmd.Stdout = &stdout
	fmt.Println(pushCmd.String())
	fmt.Println("docker push images")
	if err := pushCmd.Run(); err != nil {
		fmt.Println(stderr.String())
		return err
	}
	fmt.Println("docker push images complate")
	return nil
}

func loginAction(ctx *cli.Context) error {
	var (
		loginCmd *exec.Cmd
		stderr   bytes.Buffer
		stdout   bytes.Buffer
	)
	if needSudo() {
		loginCmd = exec.Command("sudo", "docker", "login",
			"--username=root", "--password=123456",
			"https://registry.cn-hangzhou.aliyuncs.com")
	} else {
		loginCmd = exec.Command("docker", "login", "--username=root", "--password=123456", "https://registry.cn-hangzhou.aliyuncs.com")
	}
	loginCmd.Stderr = &stderr
	loginCmd.Stdout = &stdout
	fmt.Println(loginCmd.String()) //打印执行命令
	if err := loginCmd.Run(); err != nil {
		fmt.Println(stderr.String())
		return fmt.Errorf("docker login fail, error: %s.stderr: %s", err, stderr.String())
	}
	fmt.Println("docker login success")

	return nil
}
