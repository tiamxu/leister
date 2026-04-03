package jenkins

import (
	"fmt"

	"github.com/tiamxu/kit/cli"
	"github.com/tiamxu/kit/log"
	"github.com/tiamxu/leister/client"
)

// Tool Jenkins 工具
type Tool struct {
	Client *client.Client
}

// Name 工具名称
func (t *Tool) Name() string { return "jenkins" }

// Description 工具描述
func (t *Tool) Description() string { return "Manage jenkins cmd" }

// Flags 工具全局标志
func (t *Tool) Flags() []cli.Flag {
	// 不在根命令添加 Flags，只在子命令中添加，避免解析顺序问题
	return []cli.Flag{}
}

// Commands 工具命令
func (t *Tool) Commands() []*cli.Command {
	return []*cli.Command{
		cli.NewCommand("create").
			SetDescription("Create one jenkins job").
			AddFlags(cli.StringFlag("name", "n", "", "Set jenkins appName")).
			AddFlags(cli.StringFlag("group", "g", "", "Set jenkins appGroup")).
			SetRun(func(ctx *cli.Context) error {
				return t.RunCreateJob(ctx)
			}),
		cli.NewCommand("cts").
			SetDescription("Create many jenkins jobs").
			AddFlags(cli.StringFlag("group", "g", "", "Set jenkins appGroup")).
			SetRun(func(ctx *cli.Context) error {
				return t.RunCreateJobs(ctx)
			}),
		cli.NewCommand("update").
			SetDescription("Update many jenkins jobs config").
			AddFlags(cli.StringFlag("name", "n", "", "Set jenkins appName")).
			AddFlags(cli.StringFlag("group", "g", "", "Set jenkins appGroup")).
			SetRun(func(ctx *cli.Context) error {
				return t.RunUpdateJob(ctx)
			}),
	}
}

// RunCreateJob 执行创建 Jenkins 任务命令
func (t *Tool) RunCreateJob(ctx *cli.Context) error {
	name := ctx.String("name")
	group := ctx.String("group")

	if name == "" {
		return fmt.Errorf("appName is required")
	}
	if group == "" {
		return fmt.Errorf("appGroup is required")
	}

	// 创建 Jenkins 任务请求
	req := &client.JenkinsJobRequest{
		Name:  name,
		Group: group,
	}

	// 调用 API 客户端
	resp, err := t.Client.CreateJenkinsJob(ctx.Context(), req)
	if err != nil {
		return fmt.Errorf("create jenkins job failed: %v", err)
	}

	log.Infof("Jenkins job created: %s", resp.Message)
	return nil
}

// RunCreateJobs 执行批量创建 Jenkins 任务命令
func (t *Tool) RunCreateJobs(ctx *cli.Context) error {
	group := ctx.String("group")

	if group == "" {
		return fmt.Errorf("appGroup is required")
	}

	// 这里简化处理，实际应该从数据库或其他来源获取项目列表
	// 这里假设我们有一个项目列表
	projects := []*client.JenkinsJobRequest{
		{Name: "project1", Group: group},
		{Name: "project2", Group: group},
		{Name: "project3", Group: group},
	}

	// 调用 API 客户端
	resp, err := t.Client.CreateJenkinsJobs(ctx.Context(), projects)
	if err != nil {
		return fmt.Errorf("create jenkins jobs failed: %v", err)
	}

	log.Infof("Jenkins jobs created: %s", resp.Message)
	return nil
}

// RunUpdateJob 执行更新 Jenkins 任务命令
func (t *Tool) RunUpdateJob(ctx *cli.Context) error {
	name := ctx.String("name")
	group := ctx.String("group")

	if name == "" {
		return fmt.Errorf("appName is required")
	}
	if group == "" {
		return fmt.Errorf("appGroup is required")
	}

	// 创建 Jenkins 任务请求
	req := &client.JenkinsJobRequest{
		Name:  name,
		Group: group,
	}

	// 调用 API 客户端
	resp, err := t.Client.UpdateJenkinsJob(ctx.Context(), req)
	if err != nil {
		return fmt.Errorf("update jenkins job failed: %v", err)
	}

	log.Infof("Jenkins job updated: %s", resp.Message)
	return nil
}
