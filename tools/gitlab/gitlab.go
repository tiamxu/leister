package gitlab

import (
	"errors"
	"fmt"
	"log"

	"github.com/tiamxu/kit/cli"
	"github.com/tiamxu/leister/config"
	"github.com/tiamxu/leister/database"
	"github.com/xanzy/go-gitlab"
)

type Tool struct{}

func (t *Tool) Name() string        { return "gitlab" }
func (t *Tool) Description() string { return "Manage gitlab cmd" }

func (t *Tool) Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag("name", "n", "", "Set gitlab project name"),
		cli.StringFlag("group", "g", "", "Set gitlab group"),
	}
}

func (t *Tool) Commands() []*cli.Command {
	return []*cli.Command{
		cli.NewCommand("get").
			SetDescription("Get gitlab project info console").
			AddFlags(cli.RequiredFlag(cli.StringFlag("name", "n", "", "Set gitlab project name"))).
			AddFlags(cli.RequiredFlag(cli.StringFlag("group", "g", "", "Set gitlab group"))).
			SetRun(func(ctx *cli.Context) error {
				return RunGetProject(ctx)
			}),
		cli.NewCommand("gen").
			SetDescription("Generate gitlab project data to db").
			AddFlags(cli.RequiredFlag(cli.StringFlag("group", "g", "", "Set gitlab group"))).
			SetRun(func(ctx *cli.Context) error {
				return RunGenProject(ctx)
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

func Connect(cfg *config.Config) (*gitlab.Client, error) {
	return gitlab.NewClient(cfg.Gitlab.Token, gitlab.WithBaseURL(cfg.Gitlab.Url))
}

func RunGetProject(ctx *cli.Context) error {
	return getProject(ctx)
}

func RunGenProject(ctx *cli.Context) error {
	return genProjects(ctx)
}

func genProjects(ctx *cli.Context) error {
	appGroup := ctx.String("group")
	if appGroup == "" {
		return errors.New("required OPTIONS --group or -g")
	}
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

func getProject(ctx *cli.Context) error {
	appName := ctx.String("name")
	appGroup := ctx.String("group")
	if appName == "" || appGroup == "" {
		return errors.New("required OPTIONS --name(or -n) and --group(or -g) ")
	}
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
