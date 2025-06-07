package controller

import (
	"fmt"
	"github.com/goeoeo/gitx/model"
	"github.com/goeoeo/gitx/repo"
	"github.com/goeoeo/gitx/util"
	"github.com/sirupsen/logrus"
	"sort"
)

type JiraController struct {
	config *repo.Config
	jm     *model.JiraMgr
}

func NewJiraController(config *repo.Config) (jc *JiraController, err error) {
	jc = &JiraController{
		config: config,
	}

	//载入jira数据
	jc.jm, err = model.NewJiraMgr()

	return
}

// Clear 检查目标分支中是否已合并定时任务若，若已合并则标记，同时删除本地以及远程分支
// 配合定时任务，保持本地项目和远程项目的分支简洁性
func (jc *JiraController) Clear() (err error) {
	//载入jira数据
	for _, j := range jc.jm.JiraList {
		merged := 0
		for _, jb := range j.BranchList {
			if err = jc.delBranch(j, jb); err != nil {
				return
			}
			if jb.Merged {
				merged++
			}
		}

		//已全部merge
		if len(j.BranchList) == merged && merged == len(j.TargetBranch) {
			j.Merged = true
		}
	}

	//持久化
	err = jc.jm.Save()
	return
}

// Add 添加Jira
func (jc *JiraController) Add(project, jiraID string, targetBranch []string) (err error) {
	repoCfg := jc.config.Repo[project]
	if repoCfg == nil {
		return fmt.Errorf("项目仓库信息缺失:%s", project)
	}

	if len(targetBranch) == 0 {
		return fmt.Errorf("目标分支不能为空")
	}

	return jc.jm.AddJira(project, jiraID, targetBranch)
}

// Del 删除Jira
func (jc *JiraController) Del(project, jiraID string) (err error) {
	return jc.jm.DelJira(project, jiraID)
}

func (jc *JiraController) delBranch(j *model.Jira, jb *model.JiraBranch) (err error) {
	var (
		merged bool
	)
	if jb.DevBranch == "" || jb.Merged {
		return
	}

	repoCfg := jc.config.Repo[j.Project]
	if repoCfg == nil {
		return fmt.Errorf("项目仓库信息缺失:%s", j.Project)
	}

	git := repo.NewGitRepo(repoCfg.Path, repoCfg.Url)

	if merged, err = jc.checkBranchMerged(j, jb); err != nil {
		return
	}
	if !merged {
		logrus.Infof("跳过，分支未合并:%s \n", jb.BranchName)
		return
	}

	//包含后标记
	jb.Merged = true
	logrus.Infof("正在删除分支:%s \n", jb.BranchName)

	//删除远程分支
	if err := git.DelRemoteBranch(jb.BranchName); err != nil {
		logrus.Debugf("删除远程分支错误:%s\n", err)
	}

	//删除本地分支
	if err := git.DelLocalBranch(jb.BranchName); err != nil {
		logrus.Debugf("删除远程分支错误:%s\n", err)
	}

	return
}

// checkBranchMerged 检查分支是否合并
func (jc *JiraController) checkBranchMerged(j *model.Jira, jb *model.JiraBranch) (merged bool, err error) {
	var (
		commits []*model.CommitInfo
	)

	if jb.BranchName == "" || jb.Merged {
		return
	}
	repoCfg := jc.config.Repo[j.Project]
	if repoCfg == nil {
		return false, fmt.Errorf("项目仓库信息缺失:%s", j.Project)
	}

	git := repo.NewGitRepo(repoCfg.Path, repoCfg.Url)
	//check 目标分支
	if err = git.SwitchBranch(jb.TargetBranch); err != nil {
		return
	}
	defer func() {
		_ = git.ResetBranch()
	}()

	//拉取最新代码
	if err = git.Pull(); err != nil {
		if err = git.ResetBranch(); err != nil {
			return
		}

		if err = git.DelLocalBranch(jb.TargetBranch); err != nil {
			return
		}

		if err = git.NewBranchFromRemote(jb.TargetBranch); err != nil {
			return
		}

		if err = git.Pull(); err != nil {
			return
		}

	}

	if commits, err = git.GetCommitInfo(j.GetCherryPickMsg()); err != nil {
		return
	}

	lci := jb.LastCommitInfo()
	if lci == nil {
		return
	}

	after := false
	for _, c := range commits {
		if c.CreateTime.After(lci.CreateTime) {
			after = true
			break
		}
	}

	//远程分支的提交信息没有比 jb中最大的时间大，说明远程还未合入
	if !after {
		return
	}

	//已合并
	return true, nil

}

// Print 打印出那些为合并完成的Jira
func (jc *JiraController) Print(jiraId string) {
	var (
		rows [][]string
	)

	if err := jc.syncMergeInfo(); err != nil {
		logrus.Infof("同步merge信息错误:%v", err)
		return
	}

	for _, jr := range jc.jm.JiraList {
		if jr.Complete() {
			continue
		}

		if jiraId != "" && jr.JiraID != jiraId {
			continue
		}
		sort.Slice(jr.BranchList, func(i, j int) bool {
			return jr.BranchList[i].TargetBranch < jr.BranchList[j].TargetBranch
		})

		rows = append(rows, []string{jr.GetDesc(), "MR", "状态", "更新时间"})

		for _, jb := range jr.BranchList {
			status := "待提交"
			if jb.DevBranch != "" {
				status = "待合并"
			}

			if jb.DevBranch != "" && jb.Merged {
				status = "已合并"
			}

			rows = append(rows, []string{fmt.Sprintf("%s=>%s", jb.DevBranch, jb.TargetBranch), jb.MR(), status, jb.UpdateTime.Format("2006-01-02 15-04-05")})
		}

		l := ""
		rows = append(rows, []string{l, l, l})
	}
	if len(rows) > 0 {
		rows = rows[0 : len(rows)-1]
	}

	util.PrintTable(rows, nil)
}

// syncMergeInfo 合并同步信息
func (jc *JiraController) syncMergeInfo() (err error) {
	var (
		merged   bool
		saveData bool
	)
	for _, jr := range jc.jm.JiraList {
		for _, jb := range jr.BranchList {
			if jb.Merged || jb.DevBranch == "" {
				continue
			}

			if merged, err = jc.checkBranchMerged(jr, jb); err != nil {
				return
			}

			if merged {
				saveData = true
				jb.Merged = true
			}
		}
	}

	if !saveData {
		return
	}

	err = jc.jm.Save()
	return
}
