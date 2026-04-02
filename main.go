package main

import (
	"github.com/tiamxu/kit/cli"
	"github.com/tiamxu/kit/log"
	"github.com/tiamxu/leister/config"
	"github.com/tiamxu/leister/database"
	"github.com/tiamxu/leister/tools/docker"
	"github.com/tiamxu/leister/tools/gitlab"
	"github.com/tiamxu/leister/tools/jenkins"
	"github.com/tiamxu/leister/tools/kube"
)

func main() {
	// 初始化日志
	if err := log.InitLogger(&log.Config{
		Level:  "info",
		Type:   "stdout",
		Format: "text",
	}); err != nil {
		panic(err)
	}
	defer log.Sync()

	// 加载配置
	cfg := config.Load()

	// 初始化数据库连接
	if err := database.Connect(&cfg.DB); err != nil {
		log.Infoln("Database connection failed: %v", err)
	}

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
		&docker.Tool{},
		&gitlab.Tool{},
		&jenkins.Tool{},
		&kube.Tool{},
	)

	if err := app.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
