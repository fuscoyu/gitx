package cmd

import (
	"github.com/goeoeo/gitx/repo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var HookCmd = &cobra.Command{
	Use:   "hook",
	Short: "通过Hook可以触发jenkins刷代码",
	Run: func(cmd *cobra.Command, args []string) {
		config := repo.GetConfig(configPath).Init()

		if project == "" {
			project = config.Patch.CurrentProject
		}

		r, err := config.CurrentRepo()

		config.CheckErr(err)

		if branchList == "" {
			logrus.Fatal("分支名不能为空")
		}

		p := repo.NewRepoPush(r, config, branchList, nil, false)
		p.AutoMergeBranchHook()

	},
}

func init() {
	HookCmd.Flags().StringVarP(&project, "project", "p", "", "项目")
	HookCmd.Flags().StringVarP(&branchList, "branch", "b", "", "目标分支")
}
