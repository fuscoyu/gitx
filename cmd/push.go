package cmd

import (
	"fmt"
	"github.com/goeoeo/gitx/repo"
	"github.com/goeoeo/gitx/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var (
	PushCmd = &cobra.Command{
		Use:   "push",
		Short: "cherry-pick方式推送patch",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err                     error
				mergeUrls, tmpMergeUrls []*repo.RepoPushResult
			)
			config := repo.GetConfig(configPath)
			if debug {
				config.LogLevel = 5
			}
			config.Patch.AutoMergeHook = true
			if disableAutoMergeHook {
				config.Patch.AutoMergeHook = false
			}

			config.Init().ParseJIRA(jiraID)

			if branchList != "" {
				config.Patch.TgtBranchs = strings.Split(branchList, ",")
			}

			if planTgtBranchList != "" {
				config.Patch.PlanTgtBranchList = strings.Split(planTgtBranchList, ",")
			}

			if project == "" {
				project = config.Patch.CurrentProject
			}

			if len(config.Patch.TgtBranchs) == 0 {
				logrus.Fatalf("目标分支不能为空")
			}

			if config.Patch.DevBranch == "" {
				logrus.Fatalf("开发分支不能为空")
			}

			if project == "" {
				logrus.Fatalf("项目不能为空")
			}

			for _, project := range strings.Split(project, ",") {
				if project == "" {
					continue
				}
				tmpMergeUrls, err = pushProject(project, config)
				config.CheckErr(err)

				mergeUrls = append(mergeUrls, tmpMergeUrls...)
			}

			if len(mergeUrls) == 0 {
				return
			}

			logrus.Debugf("patch push ok! \n\n")

			fmt.Println("result:")
			var rows [][]string
			for _, row := range mergeUrls {
				desc := ""
				if len(row.OutCommits) > 0 {
					desc = row.OutCommits[0].Desc
				}
				desc = strings.Replace(desc, " ", "", -1)
				rows = append(rows, []string{row.Project, fmt.Sprintf("%s=>%s", row.DevBranch, row.TargetBranch), desc, row.MergeUrl, row.MergeRes})
			}
			//汇总打印
			util.PrintTable(rows, []string{"项目", "分支", "描述", "MR", "已合入"})
		},
	}
)

func init() {
	PushCmd.Flags().StringVarP(&configPath, "config", "c", defaultConfigPath(), "配置文件路径")
	PushCmd.Flags().StringVarP(&project, "project", "p", "", "项目")
	PushCmd.Flags().StringVarP(&jiraID, "jiraId", "j", "", "jiraID")
	PushCmd.Flags().StringVarP(&branchList, "branchList", "b", "", "目标分支，支持逗号分隔")
	PushCmd.Flags().StringVarP(&planTgtBranchList, "planTgtBranchList", "t", "", "计划要推的分支列表,逗号分隔")
	PushCmd.PersistentFlags().BoolVarP(&disableAutoMergeHook, "disableAutoMergeHook", "a", false, "自动合并不执行hook")
	PushCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "忽略本地记录，cherry-pick所有commit")
	PushCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "开启debug日志")
}

func pushProject(project string, config *repo.Config) (mergeUrls []*repo.RepoPushResult, err error) {
	r, ok := config.Repo[project]
	if !ok {
		logrus.Debugf("找不到项目仓库信息:%s\n", project)
		return
	}

	repoPatch := repo.NewRepoPatch(r, config).IgnoreLocalCommit(force)
	mergeUrls, err = repoPatch.Push()
	if err != nil {
		logrus.Debugf("git repo patch repo faild: repo: %s, err: %v \n", r.Path, err)
		return
	}
	return
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Err:%s", err)
	}
}
