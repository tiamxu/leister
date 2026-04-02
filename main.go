package main

import (
	"github.com/tiamxu/kit/cli"
	"github.com/tiamxu/kit/log"
	"github.com/tiamxu/leister/client"
	"github.com/tiamxu/leister/config"
	"github.com/tiamxu/leister/tools/docker"
	"github.com/tiamxu/leister/tools/gitlab"
	"github.com/tiamxu/leister/tools/jenkins"
	"github.com/tiamxu/leister/tools/kube"
)

func main() {
	// 初始化日志（使用默认值）
	if err := log.InitLogger(&log.Config{
		Level:  "info",
		Type:   "stdout",
		Format: "text",
	}); err != nil {
		panic(err)
	}
	defer log.Sync()

	// 加载配置（从环境变量或默认值）
	cfg := config.Load()
	log.Infof("API URL: %s, Timeout: %d", cfg.API.BaseURL, cfg.API.Timeout)

	// 初始化 API 客户端
	apiClient := client.NewClient(cfg)

	// 初始化应用
	app := cli.NewApp(cli.AppConfig{
		Name:        "gigctl",
		Description: "DevOps tools for building, CI/CD and deployment",
		Version:     "0.0.1",
	})

	// 注册工具
	app.RegisterTool(
		&docker.Tool{},
		&kube.Tool{},
		&jenkins.Tool{Client: apiClient},
		&gitlab.Tool{Client: apiClient},
	)

	// 启动应用
	if err := app.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
