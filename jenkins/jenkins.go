package jenkins

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tiamxu/kit/log"
	"github.com/tiamxu/leister/client"
	"github.com/tiamxu/leister/types"
)

// Tool Jenkins 工具
type Tool struct {
	Client *client.Client
}

// NewTool 创建 Jenkins 工具实例
func NewTool(c *client.Client) *Tool { return &Tool{Client: c} }

// AddCommands 注册 Jenkins 命令到根命令
func (t *Tool) AddCommands(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "jks",
		Short: "Manage jenkins cmd",
	}
	cmd.AddCommand(t.createCmd(), t.ctsCmd(), t.updateCmd())
	root.AddCommand(cmd)
}

func (t *Tool) createCmd() *cobra.Command {
	var name, group string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create one jenkins job",
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.RunCreateJob(cmd.Context(), name, group)
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Set jenkins appName")
	cmd.Flags().StringVarP(&group, "group", "g", "", "Set jenkins appGroup")
	return cmd
}

func (t *Tool) ctsCmd() *cobra.Command {
	var group string
	cmd := &cobra.Command{
		Use:   "cts",
		Short: "Create many jenkins jobs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.RunCreateJobs(cmd.Context(), group)
		},
	}
	cmd.Flags().StringVarP(&group, "group", "g", "", "Set jenkins appGroup")
	return cmd
}

func (t *Tool) updateCmd() *cobra.Command {
	var name, group string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update many jenkins jobs config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.RunUpdateJob(cmd.Context(), name, group)
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", "Set jenkins appName")
	cmd.Flags().StringVarP(&group, "group", "g", "", "Set jenkins appGroup")
	return cmd
}

// RunCreateJob 执行创建 Jenkins 任务命令
func (t *Tool) RunCreateJob(ctx context.Context, name, group string) error {
	if name == "" {
		return fmt.Errorf("appName is required")
	}
	if group == "" {
		return fmt.Errorf("appGroup is required")
	}

	req := &types.JenkinsJobRequest{
		Name:  name,
		Group: group,
	}

	resp, err := t.Client.CreateJenkinsJob(ctx, req)
	if err != nil {
		return fmt.Errorf("create jenkins job failed: %w", err)
	}

	log.Infof("Jenkins job created: %s", resp.Message)
	return nil
}

// RunCreateJobs 执行批量创建 Jenkins 任务命令
func (t *Tool) RunCreateJobs(ctx context.Context, group string) error {
	if group == "" {
		return fmt.Errorf("appGroup is required")
	}

	// 先从 API 拉取真实项目列表（数据源：GitLab）
	listResp, err := t.Client.ListJenkinsProjects(ctx, group)
	if err != nil {
		return fmt.Errorf("list projects failed: %w", err)
	}
	if len(listResp.Projects) == 0 {
		log.Warnf("组 %s 下未找到任何项目，无需创建", group)
		return nil
	}

	projects := make([]*types.JenkinsJobRequest, 0, len(listResp.Projects))
	for _, p := range listResp.Projects {
		projects = append(projects, &types.JenkinsJobRequest{
			Name:  p.Name,
			Group: p.Group,
		})
	}

	log.Infof("准备批量创建 %d 个 Jenkins 任务", len(projects))
	resp, err := t.Client.CreateJenkinsJobs(ctx, projects)
	if err != nil {
		return fmt.Errorf("create jenkins jobs failed: %w", err)
	}

	log.Infof("Jenkins jobs created: %s", resp.Message)
	return nil
}

// RunUpdateJob 执行更新 Jenkins 任务命令
func (t *Tool) RunUpdateJob(ctx context.Context, name, group string) error {
	if name == "" {
		return fmt.Errorf("appName is required")
	}
	if group == "" {
		return fmt.Errorf("appGroup is required")
	}

	req := &types.JenkinsJobRequest{
		Name:  name,
		Group: group,
	}

	resp, err := t.Client.UpdateJenkinsJob(ctx, req)
	if err != nil {
		return fmt.Errorf("update jenkins job failed: %w", err)
	}

	log.Infof("Jenkins job updated: %s", resp.Message)
	return nil
}
