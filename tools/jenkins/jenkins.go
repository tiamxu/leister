package jenkins

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/tiamxu/kit/cli"
	"github.com/tiamxu/leister/config"
	"github.com/tiamxu/leister/database"
)

type Tool struct{}

func (t *Tool) Name() string        { return "jks" }
func (t *Tool) Description() string { return "Manage jenkins cmd" }

func (t *Tool) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag("name", "n", "", "Set jenkins appName"),
		cli.StringFlag("group", "g", "", "Set jenkins appGroup"),
	}
}

func (t *Tool) Commands() []*cli.Command {
	return []*cli.Command{
		cli.NewCommand("create").
			SetDescription("Create one jenkins job").
			AddFlags(cli.RequiredFlag(cli.StringFlag("name", "n", "", "Set jenkins appName"))).
			AddFlags(cli.RequiredFlag(cli.StringFlag("group", "g", "", "Set jenkins appGroup"))).
			SetRun(func(ctx *cli.Context) error {
				return RunCreateJob(ctx)
			}),
		cli.NewCommand("cts").
			SetDescription("Create many jenkins jobs").
			SetRun(func(ctx *cli.Context) error {
				return RunCreateJobs(ctx)
			}),
		cli.NewCommand("update").
			SetDescription("Update many jenkins jobs config").
			SetRun(func(ctx *cli.Context) error {
				return RunUpdateJobs(ctx)
			}),
	}
}

var cfg *config.Config

func init() {
	//load config
	loadConfig()
}

func loadConfig() {
	cfg = config.Load()
}

func Connect(cfg *config.Config, ctx context.Context) (*Jenkins, error) {
	// 这里需要实现 Jenkins 连接逻辑
	// 暂时返回 nil，实际使用时需要实现
	return nil, nil
}

type Jenkins struct{}

func (j *Jenkins) GetJob(ctx context.Context, name string) (*Job, error) {
	// 这里需要实现获取 Jenkins 任务的逻辑
	// 暂时返回 nil，实际使用时需要实现
	return nil, nil
}

func (j *Jenkins) CreateJob(ctx context.Context, config string, name string) (*Job, error) {
	// 这里需要实现创建 Jenkins 任务的逻辑
	// 暂时返回 nil，实际使用时需要实现
	return nil, nil
}

func (j *Jenkins) UpdateJob(ctx context.Context, name string, config string) error {
	// 这里需要实现更新 Jenkins 任务的逻辑
	// 暂时返回 nil，实际使用时需要实现
	return nil
}

type Job struct{
	name string
}

func (j *Job) GetName() string {
	return j.name
}

func RunCreateJob(ctx *cli.Context) error {
	return createJob(ctx)
}

func RunCreateJobs(ctx *cli.Context) error {
	return createJobs(ctx)
}

func RunUpdateJobs(ctx *cli.Context) error {
	return updateJobs(ctx)
}

func createJob(ctx *cli.Context) error {
	appName := ctx.String("name")
	appGroup := ctx.String("group")
	if appName == "" || appGroup == "" {
		return errors.New("required OPTIONS --name(or -n) and --group(or -g)")
	}
	ctxBg := context.Background()
	jenkins, err := Connect(cfg, ctxBg)
	if err != nil {
		return fmt.Errorf("jenkins init error:%v", err)
	}

	job, _ := jenkins.GetJob(ctxBg, appName)
	if job != nil {
		configStr := genConfigString(config.JenkinsJobConfig, appName, appGroup)
		fmt.Printf("%s 任务存在,开始更新...\n", job.GetName())
		_ = jenkins.UpdateJob(ctxBg, appName, configStr)
		fmt.Printf("%s更新成功\n", job.GetName())
	} else {
		configStr := genConfigString(config.JenkinsJobConfig, appName, appGroup)
		fmt.Printf("开始新建任务%s\n", appName)
		_, err = jenkins.CreateJob(ctxBg, configStr, appName)
		if err != nil {
			return fmt.Errorf("jenkins create job error:%v", err)
		}
		fmt.Printf("%s创建成功\n", appName)
	}

	return nil
}

func createJobs(ctx *cli.Context) error {
	appName := ctx.String("name")
	appGroup := ctx.String("group")
	initMySQL()
	ctxBg := context.Background()
	jenkins, err := Connect(cfg, ctxBg)
	if err != nil {
		return fmt.Errorf("jenkins init error:%v", err)
	}
	var items = []database.Item{}
	if appName == "" && appGroup == "" {
		items, err = database.GetAllItemData()
		if err != nil {
			return errors.New("GetAllItemData Error")
		}
	} else if appName != "" && appGroup != "" {
		items, err = database.SelectItemByWhereWithName(appName, appGroup)
		if err != nil {
			return errors.New("SelectItemByWhereWitchName Error")
		}
	}
	fmt.Printf("items:%v\n", items)

	if len(items) == 0 {
		return errors.New("NOT Found App Service")
	}
	for _, item := range items {
		appName = item.AppName
		appGroup = item.AppGroup
		job, _ := jenkins.GetJob(ctxBg, appName)
		if job != nil {
			continue
		}
		fmt.Printf("开始新建任务%s\n", appName)
		configStr := genConfigString(config.JenkinsJobConfig, appName, appGroup)
		_, err = jenkins.CreateJob(ctxBg, configStr, appName)
		if err != nil {
			return fmt.Errorf("jenkins create job error:%v", err)
		}
		fmt.Printf("%s创建成功\n", appName)
	}

	return nil
}

func updateJobs(ctx *cli.Context) error {
	appName := ctx.String("name")
	appGroup := ctx.String("group")
	initMySQL()
	ctxBg := context.Background()
	jenkins, err := Connect(cfg, ctxBg)
	if err != nil {
		return fmt.Errorf("jenkins init error:%v", err)
	}
	var items = []database.Item{}
	if appName == "" && appGroup == "" {
		items, err = database.GetAllItemData()
		if err != nil {
			return errors.New("GetAllItemData Error")
		}
	} else if appName != "" && appGroup != "" {
		items, err = database.SelectItemByWhereWithName(appName, appGroup)
		if err != nil {
			return errors.New("SelectItemByWhereWitchName Error")
		}
	}
	fmt.Printf("items:%v\n", items)

	if len(items) == 0 {
		return errors.New("NOT Found App Service")
	}
	for _, item := range items {
		appName = item.AppName
		appGroup = item.AppGroup
		configStr := genConfigString(config.JenkinsJobConfig, appName, appGroup)
		job, err := jenkins.GetJob(ctxBg, appName)
		if err != nil {
			continue
		}
		if job != nil {
			fmt.Printf("%s 任务开始更新...\n", job.GetName())
			jenkins.UpdateJob(ctxBg, appName, configStr)
			fmt.Printf("%s更新成功\n", job.GetName())
		}
	}

	return nil
}

func initMySQL() {
	dbConfig := &database.Config{
		Driver:          cfg.DB.Driver,
		Database:        cfg.DB.Database,
		Username:        cfg.DB.Username,
		Password:        cfg.DB.Password,
		Host:            cfg.DB.Host,
		Port:            cfg.DB.Port,
		MaxIdleConns:    cfg.DB.MaxIdleConns,
		MaxOpenConns:    cfg.DB.MaxOpenConns,
		ConnMaxLifetime: cfg.DB.ConnMaxLifetime,
	}
	err := database.Connect(dbConfig)
	if err != nil {
		fmt.Printf("DB connect failed error:%v\n", err)
		return
	}
}

// generate job configuration file
func genConfigString(tpl, name, group string) string {
	str := strings.Replace(tpl, "${app_name}", name, -1)
	str = strings.Replace(str, "${app_group}", group, -1)
	return str
}
