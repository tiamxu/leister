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
func genConfigString(name, group string) string {
	str := strings.Replace(config.JenkinsJobConfig, "${app_name}", name, -1)
	str = strings.Replace(str, "${app_group}", group, -1)
	return str
}
func RunCreateJob(c *cli.Context) error {
	return createJob(c)
}
func RunCreateJobs(c *cli.Context) error {
	return createJobs(c)
}
func createJob(c *cli.Context) error {
	ctx := context.Background()
	jenkins, err := Connect(cfg, ctx)
	if err != nil {
		return fmt.Errorf("jenkins init error:%v", err)
	}
	fmt.Printf("开始新建任务%s\n", appName)
	configStr := genConfigString(appName, appGroup)
	_, err = jenkins.CreateJob(ctx, configStr, appName)
	if err != nil {
		return fmt.Errorf("jenkins create job error:%v", err)
	}
	fmt.Printf("%s创建成功\n", appName)
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
		fmt.Printf("开始新建任务%s\n", appName)
		configStr := genConfigString(appName, appGroup)
		_, err = jenkins.CreateJob(ctx, configStr, appName)
		if err != nil {
			return fmt.Errorf("jenkins create job error:%v", err)
		}
		fmt.Printf("%s创建成功\n", appName)
	}

	return nil
}

func initMySQL() {
	err := database.Connect(cfg.DB)
	if err != nil {
		fmt.Printf("DB connect failed error:%v\n", err)
	}
	fmt.Println("MySQL Connect success....")
}
