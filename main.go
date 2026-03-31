package main

import (
	"github.com/tiamxu/kit/cli"
	"github.com/tiamxu/leister/tools/build"
	"github.com/tiamxu/leister/tools/jenkins"
	"github.com/tiamxu/leister/tools/gitlab"
	"github.com/tiamxu/leister/tools/kube"
)

func main() {
	app := cli.NewApp(cli.AppConfig{
		Name:        "gigctl",
		Description: "DevOps tools for building, CI/CD and deployment",
		Version:     "0.0.1",
		PreRun: func(ctx *cli.Context) error {
			// 加载配置等初始化操作
			return nil
		},
	})

	app.RegisterTool(
		&build.Tool{},
		&jenkins.Tool{},
		&gitlab.Tool{},
		&kube.Tool{},
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
