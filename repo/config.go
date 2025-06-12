package repo

import (
	"bytes"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/goeoeo/gitx/model"
	"github.com/goeoeo/gitx/util"
	"github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog/zt_formatter"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var cfg *Config

type Config struct {
	Repo          map[string]*Repo `yaml:"repo"`
	Patch         *Patch           `yaml:"patch"`
	HomeDir       string           `yaml:"home_dir"`
	LogLevel      int              `yaml:"log_level"`
	GitLabConfigs []*GitLabConfig  `yaml:"gitLab_configs"`
	pwd           string
	logBuffer     bytes.Buffer
}

type Repo struct {
	Name                string            `yaml:"name"`
	Url                 string            `yaml:"url"`
	Path                string            `yaml:"path"`
	CreateMr            bool              `yaml:"create_mr"`              //自动创建mr
	AutoMergeBranchList []string          `yaml:"auto_merge_branch_list"` //自动合并的分支
	AutoMergeBranchHook map[string]string `yaml:"auto_merge_branch_hook"` //自动合并分支后触发的操作
}

type Patch struct {
	DevBranch         string            `yaml:"dev_branch"`
	TgtBranchs        []string          `yaml:"tgt_branchs"`
	PlanTgtBranchList []string          `yaml:"plan_tgt_branch_list"` //计划要推的分支列表
	BranchAlias       map[string]string `yaml:"branch_alias"`         //分支别名
	JiraId            string            `yaml:"jira_id"`
	JiraDesc          string            `yaml:"jira_desc"`
	CommitType        string            `yaml:"commit_type"`     //提交的类型，可以是jira,也可以是整个message
	CommitMsg         string            `yaml:"commit_msg"`      //提交的message
	CurrentProject    string            `yaml:"current_project"` //当前项目，可通过pwd进行推断
	JiraProjects      []string          `yaml:"jira_projects"`   //jira项目，用于推断CommitType
	TmpBranchFmt      string            `yaml:"tmp_branch_fmt"`  //临时分支的格式默认：{jiraID}_{jiraDesc}_{tgtBranch}
	AutoMergeHook     bool              `yaml:"auto_merge_hook"` // 是否执行hook
}

type GitLabConfig struct {
	BaseUrl string `yaml:"base_url"` //https://git.internal.yunify.com
	Token   string `yaml:"token"`
}

func GetConfig(configPaths ...string) *Config {
	if cfg != nil {
		return cfg
	}

	var configPath string
	if len(configPaths) > 0 {
		configPath = configPaths[0]
	}
	if configPath == "" {
		d, _ := os.UserHomeDir()
		configPath = d + "/.patch/config.yaml"
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var config *Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	if config.Patch == nil {
		config.Patch = &Patch{}
	}

	if config.Patch.BranchAlias == nil {
		config.Patch.BranchAlias = make(map[string]string)
	}

	if config.HomeDir == "" {
		config.HomeDir, _ = os.UserHomeDir()
		config.HomeDir += "/.patch"
	}

	if config.Patch.TmpBranchFmt == "" {
		config.Patch.TmpBranchFmt = "{jiraID}_{jiraDesc}_{tgtBranch}"
	}
	cfg = config
	return cfg
}

// 解析当前目录信息
func (c *Config) parsePwd() (err error) {
	c.pwd = os.Getenv("PWD")
	logrus.Debugf("当前目录:%s", c.pwd)

	projectName := util.GetLastDir(c.pwd)

	r := &Repo{
		Name:     util.GetLastDir(c.pwd),
		Path:     c.pwd,
		CreateMr: true,
	}

	if r.Url, err = util.FindOriginURL(c.pwd); err != nil {
		return
	}

	find := false
	for _, v := range c.Repo {
		//配置文件中存在配置
		if v.Name == projectName {
			find = true
			v.Url = util.Default(v.Url, r.Url)
			v.Path = util.Default(v.Path, r.Path)
			break
		}

	}

	if !find {
		c.Repo[projectName] = r
	}

	return

}

func (c *Config) Init() *Config {
	if err := c.parsePwd(); err != nil {
		logrus.Fatalf("当前目标非git项目目录，err:%s", err)
	}

	if c.Patch.JiraDesc == "" {
		c.Patch.JiraDesc = "x"
	}

	if c.Patch.CurrentProject == "" {
		for k, v := range c.Repo {
			if v.Path == c.pwd {
				c.Patch.CurrentProject = k
			}
		}
	}

	if c.Patch.CommitType == "" {
		c.Patch.CommitType = model.CommitTypeJira
	}

	//自动解析当前分支
	if c.Patch.DevBranch == "" {
		c.Patch.DevBranch = AutoBranch(c.pwd)
	}

	c.InitLog()

	logrus.Debugf("get c ok! \n %v", c.Patch)
	return c
}

func (c *Config) GetGitLabConfig(url string) *GitLabConfig {
	for _, v := range c.GitLabConfigs {
		if strings.HasPrefix(url, v.BaseUrl) {
			return v
		}
	}
	return nil
}

func (c *Config) GetRepo(name string) *Repo {
	return c.Repo[name]
}

func (c *Config) CurrentRepo() (r *Repo, err error) {
	r, ok := c.Repo[c.Patch.CurrentProject]
	if ok {
		return r, nil
	}
	return nil, fmt.Errorf("当前repo不存在:[%s]", c.Patch.CurrentProject)
}

func (c *Config) Print() {
	content, _ := yaml.Marshal(c)
	fmt.Println(string(content))
}

func (c *Config) ParseJIRA(jiraID string) *Config {
	if jiraID != "" {
		c.Patch.JiraId = jiraID
	}

	//自动解析jiraID
	if c.Patch.JiraId == "" {
		c.Patch.JiraId, c.Patch.CommitType, c.Patch.CommitMsg = AutoJiraID(c.pwd, c.Patch.JiraProjects, "")
		return c
	}

	//通过commit挑选
	re := regexp.MustCompile(`^[a-z0-9]+$`)
	if re.MatchString(c.Patch.JiraId) {
		c.Patch.JiraId, c.Patch.CommitType, c.Patch.CommitMsg = AutoJiraID(c.pwd, c.Patch.JiraProjects, c.Patch.JiraId)
		return c
	}

	return c

}

func (c *Config) InitLog() {

	if c.LogLevel > 0 {
		logrus.SetLevel(logrus.Level(c.LogLevel))
		return
	}

	//logLevel为0的情况,日志级别调整为debug，并将输出定向到文件
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&zt_formatter.ZtFormatter{
		Formatter: nested.Formatter{
			TimestampFormat: "2006-01-02T15:04:05.000",
			HideKeys:        false,
			TrimMessages:    true,
			ShowFullLevel:   true,
		},
	})
	logrus.SetOutput(&c.logBuffer)

}

func (c *Config) CheckErr(err error) {
	if err == nil {
		return
	}

	logBytes := c.logBuffer.Bytes()
	os.Stdout.Write(logBytes)

	logrus.Fatalf("err:%s\n", err)
}

func (p *Patch) GetTgtBranchs() (res []string) {
	for _, v := range p.TgtBranchs {
		branchName := v
		alias, ok := p.BranchAlias[v]
		if ok {
			branchName = alias
		}

		res = append(res, branchName)
	}
	return
}

func (p *Patch) GetPlanTgtBranchList() (res []string) {
	for _, v := range p.PlanTgtBranchList {
		branchName := v
		alias, ok := p.BranchAlias[v]
		if ok {
			branchName = alias
		}

		res = append(res, branchName)
	}
	return
}
