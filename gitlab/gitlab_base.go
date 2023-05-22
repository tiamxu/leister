package gitlab

import (
	"errors"
	"fmt"
	"log"

	"github.com/tiamxu/leister/database"

	"github.com/urfave/cli/v2"
	"github.com/xanzy/go-gitlab"
)

var (
	appName  string
	appGroup string
	Flags    = []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "set gitlab project name",
		},
		&cli.StringFlag{
			Name:    "group",
			Aliases: []string{"g"},
			Usage:   "set gitlab  group",
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
		return errors.New("required OPTIONS --name(or -n) and --group(or -g) ")
	}
	return nil
}

func InitGenFlags(c *cli.Context) error {
	appName = c.String("name")
	appGroup = c.String("group")
	if appGroup == "" {
		return errors.New("required OPTIONS --group or -g")
	}
	initMySQL()
	return nil
}
func RunGetProject(c *cli.Context) error {
	return getProject(c)
}
func RunGenProject(c *cli.Context) error {
	return genProjects(c)
}

func genProjects(c *cli.Context) error {
	var item = database.Item{}
	var items = []database.Item{}
	git, err := Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	var gid int
	groupOption := &gitlab.ListGroupsOptions{Search: gitlab.String(appGroup)}
	groups, _, err := git.Groups.ListGroups(groupOption)
	if err != nil {
		log.Fatalf("Failed to get groups err: %v", err)
	}
	group := groups[0]
	gid = group.ID
	opt := &gitlab.ListGroupProjectsOptions{ListOptions: gitlab.ListOptions{Page: 1, PerPage: 50}}
	projects, _, err := git.Groups.ListGroupProjects(gid, opt)
	if err != nil {
		log.Fatalf("Failed to get projects err: %v", err)
	}
	for _, v := range projects {
		fmt.Println(v.ID, v.Name, v.HTTPURLToRepo, v.SSHURLToRepo)
		item.CodeID = v.ID
		item.AppName = v.Name
		item.AppGroup = appGroup
		item.AppType = "go"
		item.HTTPURLToRepo = v.HTTPURLToRepo
		item.SSHURLToRepo = v.SSHURLToRepo
		items = append(items, item)

	}

	for _, item := range items {
		n, err := database.AddItem(item)
		if err != nil {
			log.Fatalf("插入数据错误: %v", err)
		}
		fmt.Printf("insert success,affected rows%v\n", n)
	}

	return nil

}
func getProject(c *cli.Context) error {
	git, err := gitlab.NewClient(cfg.Gitlab.Token, gitlab.WithBaseURL(cfg.Gitlab.Url))
	if err != nil {
		log.Fatalf("Failed to create gitlab client: %v", err)
	}
	var gid int
	groupOption := &gitlab.ListGroupsOptions{Search: gitlab.String(appGroup)}
	groups, _, err := git.Groups.ListGroups(groupOption)
	if err != nil {
		log.Fatalf("Failed to get gitlab groups err: %v", err)
	}
	if len(groups) == 0 {
		fmt.Println("NOT Found Gitlab Group...")
		return nil
	}
	group := groups[0]
	gid = group.ID
	opt := &gitlab.ListGroupProjectsOptions{Search: gitlab.String(appName), ListOptions: gitlab.ListOptions{Page: 1, PerPage: 50}}
	projects, _, err := git.Groups.ListGroupProjects(gid, opt)
	if err != nil {
		log.Fatalf("Failed to get projects err: %v", err)
	}
	if len(projects) == 0 {
		fmt.Println("NOT Found Project...")
		return nil
	}
	for _, v := range projects {
		fmt.Printf("ProjectName: %v\n", v.Name)
		fmt.Printf("ProjectID: %v\n", v.ID)
		fmt.Printf("HTTP_URL_TO_Repo: %v\n", v.HTTPURLToRepo)
		fmt.Printf("SSH_URL_To_Repo: %v\n", v.SSHURLToRepo)
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
