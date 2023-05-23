package jenkins

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/tiamxu/leister/config"
	"github.com/tiamxu/leister/database"
	"github.com/urfave/cli/v2"
)

// 读取配置文件
var (
	appName  string
	appGroup string
	Flags    = []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "Set jenkins appName",
		},
		&cli.StringFlag{
			Name:    "group",
			Aliases: []string{"g"},
			Usage:   "Set jenkins appGroup",
		},
	}
)

func init() {
	//load config
	loadConfig()
}

func InitFlags(c *cli.Context) error {
	appName = c.String("name")
	appGroup = c.String("group")
	if appName == "" || appGroup == "" {
		return errors.New("required OPTIONS --name(or -n) and --group(or -g)")
	}
	return nil
}

func InitCtxFlags(c *cli.Context) error {
	appName = c.String("name")
	appGroup = c.String("group")
	if appName != "" && appGroup == "" {
		return errors.New("required OPTIONS  --group(or -g)")
	}
	initMySQL()
	return nil
}

// generate job configuration file
func genConfigString(tpl, name, group string) string {
	str := strings.Replace(tpl, "${app_name}", name, -1)
	str = strings.Replace(str, "${app_group}", group, -1)
	return str
}
func RunCreateJob(c *cli.Context) error {
	return createJob(c)
}
func RunCreateJobs(c *cli.Context) error {
	return createJobs(c)
}
func RunUpdateJobs(c *cli.Context) error {
	return updateJobs(c)
}
func createJob(c *cli.Context) error {
	ctx := context.Background()
	jenkins, err := Connect(cfg, ctx)
	if err != nil {
		return fmt.Errorf("jenkins init error:%v", err)
	}

	job, _ := jenkins.GetJob(ctx, appName)
	if job != nil {
		configStr := genConfigString(config.JenkinsJobConfig, appName, appGroup)
		fmt.Printf("%s 任务存在,开始更新...\n", job.GetName())
		_ = jenkins.UpdateJob(ctx, appName, configStr)
		fmt.Printf("%s更新成功\n", job.GetName())
	} else {
		configStr := genConfigString(config.JenkinsJobConfig, appName, appGroup)
		fmt.Printf("开始新建任务%s\n", appName)
		_, err = jenkins.CreateJob(ctx, configStr, appName)
		if err != nil {
			return fmt.Errorf("jenkins create job error:%v", err)
		}
		fmt.Printf("%s创建成功\n", appName)
	}

	return nil
}

func createJobs(c *cli.Context) error {
	ctx := context.Background()
	jenkins, err := Connect(cfg, ctx)
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
		job, _ := jenkins.GetJob(ctx, appName)
		if job != nil {
			continue
		}
		fmt.Printf("开始新建任务%s\n", appName)
		configStr := genConfigString(config.JenkinsJobConfig, appName, appGroup)
		_, err = jenkins.CreateJob(ctx, configStr, appName)
		if err != nil {
			return fmt.Errorf("jenkins create job error:%v", err)
		}
		fmt.Printf("%s创建成功\n", appName)
	}

	return nil
}

func updateJobs(c *cli.Context) error {
	ctx := context.Background()
	jenkins, err := Connect(cfg, ctx)
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
		job, err := jenkins.GetJob(ctx, appName)
		if err != nil {
			continue
		}
		if job != nil {
			fmt.Printf("%s 任务开始更新...\n", job.GetName())
			jenkins.UpdateJob(ctx, appName, configStr)
			// if err != nil {
			// 	return fmt.Errorf("jenkins update job error:%v", err)
			// }
			fmt.Printf("%s更新成功\n", job.GetName())
		}
	}

	return nil
}

func initMySQL() {
	err := database.Connect(cfg.DB)
	if err != nil {
		fmt.Printf("DB connect failed error:%v\n", err)
		return
	}
}
