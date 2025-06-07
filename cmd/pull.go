package cmd

import (
	"github.com/goeoeo/gitx/repo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

var PullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
		)

		config := repo.GetConfig(configPath)

		for _, project := range strings.Split(project, ",") {
			if project == "" {
				continue
			}
			err = pullProject(project, config)
			config.CheckErr(err)
		}
		return
	},
}

func init() {

	PullCmd.Flags().StringVarP(&configPath, "config", "c", defaultConfigPath(), "配置文件路径")
	PullCmd.Flags().StringVarP(&project, "project", "p", "", "项目")
}

func pullProject(project string, config *repo.Config) (err error) {
	r, ok := config.Repo[project]
	if !ok {
		logrus.Debugf("找不到项目仓库信息:%s\n", project)
		return
	}

	logrus.Debugf("git pull patch target branch ok: repo: %s \n", r.Path)

	repoPatch := repo.NewRepoPatch(r, config)

	if err = repoPatch.Pull(false); err != nil {
		logrus.Debugf("git pull patch repo faild: repo: %s, err: %v \n", r.Path, err)
		return
	}
	return
}
