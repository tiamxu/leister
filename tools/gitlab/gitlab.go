package gitlab

import (
	"fmt"

	"github.com/tiamxu/kit/cli"
	"github.com/tiamxu/kit/log"
	"github.com/tiamxu/leister/client"
)

// Tool GitLab 工具
type Tool struct {
	Client *client.Client
}

// Name 工具名称
func (t *Tool) Name() string { return "gitlab" }

// Description 工具描述
func (t *Tool) Description() string { return "Manage gitlab cmd" }

// Flags 工具全局标志
func (t *Tool) Flags() []cli.Flag {
	// 不在根命令添加 Flags，只在子命令中添加，避免解析顺序问题
	return []cli.Flag{}
}

// Commands 工具命令
func (t *Tool) Commands() []*cli.Command {
	return []*cli.Command{
		cli.NewCommand("get").
			SetDescription("Get gitlab project info console").
			AddFlags(cli.StringFlag("name", "n", "", "Set gitlab project name")).
			AddFlags(cli.StringFlag("group", "g", "", "Set gitlab group")).
			SetRun(func(ctx *cli.Context) error {
				return t.RunGetProject(ctx)
			}),
		cli.NewCommand("gen").
			SetDescription("Generate gitlab project data to db").
			AddFlags(cli.StringFlag("group", "g", "", "Set gitlab group")).
			SetRun(func(ctx *cli.Context) error {
				return t.RunGenProject(ctx)
			}),
	}
}

// RunGetProject 执行获取 GitLab 项目命令
func (t *Tool) RunGetProject(ctx *cli.Context) error {
	name := ctx.String("name")
	group := ctx.String("group")

	if name == "" {
		return fmt.Errorf("project name is required")
	}
	if group == "" {
		return fmt.Errorf("group is required")
	}

	// 创建 GitLab 项目请求
	req := &client.GitlabProjectRequest{
		Name:  name,
		Group: group,
	}

	// 调用 API 客户端
	resp, err := t.Client.GetGitlabProject(ctx.Context(), req)
	if err != nil {
		return fmt.Errorf("get gitlab project failed: %v", err)
	}

	// 打印项目信息
	log.Infof("GitLab Project Info:")
	log.Infof("ID: %d", resp.Project.ID)
	log.Infof("Name: %s", resp.Project.Name)
	log.Infof("Group: %s", resp.Project.Group)
	log.Infof("HTTP URL: %s", resp.Project.HTTPURLToRepo)
	log.Infof("SSH URL: %s", resp.Project.SSHURLToRepo)

	return nil
}

// RunGenProject 执行生成 GitLab 项目数据命令
func (t *Tool) RunGenProject(ctx *cli.Context) error {
	group := ctx.String("group")

	if group == "" {
		return fmt.Errorf("group is required")
	}

	// 创建 GitLab 项目生成请求
	req := &client.GitlabGenRequest{
		Group: group,
	}

	// 调用 API 客户端
	resp, err := t.Client.GenGitlabProjects(ctx.Context(), req)
	if err != nil {
		return fmt.Errorf("gen gitlab project data failed: %v", err)
	}

	// 打印生成结果
	log.Infof("GitLab Project Generation Result:")
	log.Infof("Status: %s", resp.Status)
	log.Infof("Message: %s", resp.Message)
	log.Infof("Projects generated: %d", len(resp.Projects))

	// 打印项目列表
	for _, project := range resp.Projects {
		log.Infof("- %s (ID: %d)", project.Name, project.ID)
		log.Infof("  HTTP URL: %s", project.HTTPURLToRepo)
		log.Infof("  SSH URL: %s", project.SSHURLToRepo)
	}

	return nil
}
