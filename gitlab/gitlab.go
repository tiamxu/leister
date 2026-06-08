package gitlab

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tiamxu/kit/log"
	"github.com/tiamxu/leister/client"
	"github.com/tiamxu/leister/types"
)

// Tool GitLab 工具
type Tool struct {
	Client *client.Client
}

// NewTool 创建 GitLab 工具实例
func NewTool(c *client.Client) *Tool { return &Tool{Client: c} }

// AddCommands 注册 GitLab 命令到根命令
func (t *Tool) AddCommands(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "git",
		Short: "Manage gitlab cmd",
	}
	cmd.AddCommand(t.getCmd(), t.genCmd())
	root.AddCommand(cmd)
}

func (t *Tool) getCmd() *cobra.Command {
	var name, group string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get gitlab project info console",
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.RunGetProject(cmd.Context(), name, group)
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Set gitlab project name")
	cmd.Flags().StringVarP(&group, "group", "g", "", "Set gitlab group")
	return cmd
}

func (t *Tool) genCmd() *cobra.Command {
	var group string
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate gitlab project data to db",
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.RunGenProject(cmd.Context(), group)
		},
	}
	cmd.Flags().StringVarP(&group, "group", "g", "", "Set gitlab group")
	return cmd
}

// RunGetProject 执行获取 GitLab 项目命令
func (t *Tool) RunGetProject(ctx context.Context, name, group string) error {
	if name == "" {
		return fmt.Errorf("project name is required")
	}
	if group == "" {
		return fmt.Errorf("group is required")
	}

	req := &types.GitlabProjectRequest{
		Name:  name,
		Group: group,
	}

	resp, err := t.Client.GetGitlabProject(ctx, req)
	if err != nil {
		return fmt.Errorf("get gitlab project failed: %w", err)
	}
	if resp.Project == nil {
		return fmt.Errorf("gitlab project not found: %s/%s", group, name)
	}

	log.Infof("GitLab Project Info:")
	log.Infof("ID: %d", resp.Project.ID)
	log.Infof("Name: %s", resp.Project.Name)
	log.Infof("Group: %s", resp.Project.Group)
	log.Infof("HTTP URL: %s", resp.Project.HTTPURLToRepo)
	log.Infof("SSH URL: %s", resp.Project.SSHURLToRepo)

	return nil
}

// RunGenProject 执行生成 GitLab 项目数据命令
func (t *Tool) RunGenProject(ctx context.Context, group string) error {
	if group == "" {
		return fmt.Errorf("group is required")
	}

	req := &types.GitlabGenRequest{
		Group: group,
	}

	resp, err := t.Client.GenGitlabProjects(ctx, req)
	if err != nil {
		return fmt.Errorf("gen gitlab project data failed: %w", err)
	}

	log.Infof("GitLab Project Generation Result:")
	log.Infof("Status: %s", resp.Status)
	log.Infof("Message: %s", resp.Message)
	log.Infof("Projects generated: %d", len(resp.Projects))

	for _, project := range resp.Projects {
		log.Infof("- %s (ID: %d)", project.Name, project.ID)
		log.Infof("  HTTP URL: %s", project.HTTPURLToRepo)
		log.Infof("  SSH URL: %s", project.SSHURLToRepo)
	}

	return nil
}
