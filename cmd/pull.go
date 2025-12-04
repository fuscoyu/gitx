package cmd

import (
	"strings"

	"github.com/goeoeo/gitx/repo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var PullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
		)

		config := repo.GetConfig(configPath).Init()

		// 设置目标分支：如果指定了 -b 则切换到指定分支拉取，否则直接拉取当前分支
		if branch != "" {
			config.Patch.TgtBranchs = []string{branch}
		}

		for _, project := range strings.Split(project, ",") {
			if project == "" {
				continue
			}
			err = pullProject(project, config, branch != "")
			config.CheckErr(err)
		}
	},
}

func init() {
	PullCmd.Flags().StringVarP(&configPath, "config", "c", defaultConfigPath(), "配置文件路径")
	PullCmd.Flags().StringVarP(&project, "project", "p", "", "项目")
	PullCmd.Flags().StringVarP(&branch, "branch", "b", "", "目标分支，不指定则直接拉取当前分支")
}

func pullProject(project string, config *repo.Config, switchBranch bool) (err error) {
	r, ok := config.Repo[project]
	if !ok {
		logrus.Debugf("找不到项目仓库信息:%s\n", project)
		return
	}

	if switchBranch {
		// 指定了分支，切换到目标分支并拉取
		logrus.Debugf("git pull switch to branch: repo: %s \n", r.Path)
		repoPatch := repo.NewRepoPatch(r, config)
		if err = repoPatch.Pull(false); err != nil {
			logrus.Debugf("git pull patch repo faild: repo: %s, err: %v \n", r.Path, err)
			return
		}
	} else {
		// 未指定分支，直接在当前分支拉取
		logrus.Debugf("git pull current branch: repo: %s \n", r.Path)
		gitRepo := repo.NewGitRepo(r.Path, r.Url)
		if err = gitRepo.Pull(); err != nil {
			logrus.Debugf("git pull faild: repo: %s, err: %v \n", r.Path, err)
			return
		}
	}
	return
}
