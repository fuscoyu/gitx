package model

import (
	"github.com/goeoeo/gitx/util"
	"sort"
	"strings"
	"time"
)

const (
	CommitTypeJira = "jira"
	CommitTypeMsg  = "message"
)

type (
	Jira struct {
		Project       string
		JiraID        string
		CommitType    string   //提交类型 jira=代表一个jira任务,message=代表整个commit msg
		CommitMessage string   //如果CommitType类型为message，将用此抽取 commit
		TargetBranch  []string //要合入的branch
		CreateTime    time.Time
		UpdateTime    time.Time
		BranchList    []*JiraBranch //一个jira和一个分支对应
		Merged        bool          //当所有分支合入后，标识为true，当重新patch时，更新为false
	}

	// JiraBranch JiraId 对应的branch分支
	JiraBranch struct {
		BranchName    string    ///通过cherry-pick 产生的分支名
		DevBranch     string    //分支
		TargetBranch  string    //目标分支
		Merged        bool      //是否已合入目标分支
		UpdateTime    time.Time // 更新时间
		CreateTime    time.Time
		Commits       []*CommitInfo //相关的commits
		MergeRequests []*MrInfo
		LinkInfo      *LinkInfoItem
	}

	CommitInfo struct {
		CommitId   string
		Desc       string
		CreateTime time.Time
	}
	MrInfo struct {
		Title  string
		MrId   int
		WebUrl string
	}
	//LinkInfoItem 链接发布单信息
	LinkInfoItem struct {
		LinkType string
		IssueId  string
		Summary  string
		Status   string
	}
)

func (j *Jira) Init() {
	//根据targetBranch 生成JiraBranch
	for _, branch := range j.TargetBranch {
		if j.get(branch) != nil {
			continue
		}

		j.BranchList = append(j.BranchList, &JiraBranch{
			TargetBranch: branch,
		})

	}
}

func (j *Jira) AttachBranch(branch string) *Jira {
	for _, v := range j.TargetBranch {
		if v == branch {
			return j
		}
	}

	j.TargetBranch = append(j.TargetBranch, branch)
	sort.Strings(j.TargetBranch)
	return j

}

func (j *Jira) Append(jb *JiraBranch) *Jira {
	j.UpdateTime = time.Now()
	oldJb := j.get(jb.TargetBranch)
	if oldJb != nil {
		oldJb.BranchName = jb.BranchName
		oldJb.TargetBranch = jb.TargetBranch
		oldJb.DevBranch = jb.DevBranch
		oldJb.UpdateTime = time.Now()
		oldJb.Merged = false
		oldJb.Commits = append(oldJb.Commits, jb.Commits...)
		oldJb.MergeRequests = append(oldJb.MergeRequests, jb.MergeRequests...)
		sort.SliceStable(oldJb.Commits, func(i, j int) bool {
			return oldJb.Commits[i].CreateTime.Before(oldJb.Commits[j].CreateTime)
		})

		return j
	}

	jb.CreateTime = time.Now()
	jb.UpdateTime = time.Now()

	sort.SliceStable(jb.Commits, func(i, j int) bool {
		return jb.Commits[i].CreateTime.Before(jb.Commits[j].CreateTime)
	})
	j.BranchList = append(j.BranchList, jb)

	return j
}

func (j *Jira) get(targetBranch string) *JiraBranch {
	for _, v := range j.BranchList {
		if v.TargetBranch == targetBranch {
			return v
		}
	}
	return nil
}

// BranchContainCommit 检查当前分支是否已经包含commitId
func (j *Jira) BranchContainCommit(branch, commitId string) bool {
	jb := j.get(branch)
	if jb == nil {
		return false
	}

	for _, v := range jb.Commits {
		if v.CommitId == commitId {
			return true
		}
	}

	return false
}

// Complete 判定当前jiraId 是否已完成
func (j *Jira) Complete() bool {
	for _, v := range j.BranchList {
		if !v.Merged || v.DevBranch == "" {
			return false
		}
	}

	if len(j.TargetBranch) != len(j.BranchList) {
		return false
	}

	return true
}

func (jb *JiraBranch) Desc(first bool) string {
	var desc []string
	for _, v := range jb.Commits {
		if first {
			return v.Desc
		}
		desc = append(desc, v.Desc)
	}

	desc = util.Unique(desc)
	return strings.Join(desc, ",")
}

func (jb *JiraBranch) LastCommitInfo() *CommitInfo {
	if len(jb.Commits) == 0 {
		return nil
	}

	return jb.Commits[len(jb.Commits)-1:][0]
}

func (jb *JiraBranch) MR() string {
	var arr []string
	for _, v := range jb.MergeRequests {
		arr = append(arr, v.WebUrl)
	}
	return strings.Join(arr, "\n")
}

func (j *Jira) GetCherryPickMsg() string {
	if j.CommitType == CommitTypeMsg {
		return j.CommitMessage
	}

	return j.JiraID
}

func (j *Jira) AddTargetBranch(tgt []string) *Jira {
	j.TargetBranch = append(j.TargetBranch, tgt...)
	j.TargetBranch = util.Unique(j.TargetBranch)
	return j
}

func (j *Jira) GetDesc() string {
	for _, jb := range j.BranchList {
		desc := strings.Replace(jb.Desc(false), " ", "", -1)
		if desc != "" {
			return desc
		}
	}

	return j.JiraID
}
