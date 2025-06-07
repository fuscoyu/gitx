package repo

import (
	"fmt"
	"github.com/goeoeo/gitx/util"
	"github.com/stretchr/testify/assert"
	"github.com/xanzy/go-gitlab"
	"regexp"
	"testing"
)

func TestGitRepo_GetCommitInfo(t *testing.T) {
	git := NewGitRepo("/Users/yu/code/git.internal.yunify.com/bi/pitrix-wh-daemon", "")

	ci, err := git.GetCommitInfo("BILLING-3113")
	assert.Nil(t, err)
	for _, v := range ci {
		fmt.Printf("%s,%s,%s\n", v.CommitId, v.Desc, v.CreateTime)
	}
}

func TestGirRepo_CreateMergeRequest(t *testing.T) {
	config := GetConfig().GetGitLabConfig("https://git.internal.yunify.com")

	git, err := gitlab.NewClient(config.Token, gitlab.WithBaseURL(config.BaseUrl))
	assert.Nil(t, err)

	mr, _, err := git.MergeRequests.CreateMergeRequest(3823, &gitlab.CreateMergeRequestOptions{
		Title:                stringPtr("VM-1888:1"),
		Description:          stringPtr("xx"),
		SourceBranch:         stringPtr("VM-1888_x_dev"),
		TargetBranch:         stringPtr("dev"),
		Labels:               nil,
		AssigneeID:           nil,
		AssigneeIDs:          nil,
		ReviewerIDs:          nil,
		TargetProjectID:      nil,
		MilestoneID:          nil,
		RemoveSourceBranch:   boolPtr(true),
		Squash:               boolPtr(true),
		AllowCollaboration:   nil,
		ApprovalsBeforeMerge: nil,
	})
	assert.Nil(t, err)
	//gitlab.MergeRequest{ID:44073, IID:18
	fmt.Println(mr)

}

func TestGirRepo_MergeMergeRequest(t *testing.T) {
	config := GetConfig().GetGitLabConfig("https://git.yunify.com/")

	git, err := gitlab.NewClient(config.Token, gitlab.WithBaseURL(config.BaseUrl))
	assert.Nil(t, err)

	mr, _, err := git.MergeRequests.AcceptMergeRequest(GetConfig().GetRepo("dev-tool").Url, 18, nil)
	assert.Nil(t, err)
	fmt.Println(mr)

}

func TestGirRepo_GetMergeRequest(t *testing.T) {
	config := GetConfig().GetGitLabConfig("https://git.yunify.com/")

	git, err := gitlab.NewClient(config.Token, gitlab.WithBaseURL(config.BaseUrl))
	assert.Nil(t, err)

	mr, _, err := git.MergeRequests.GetMergeRequest(GetConfig().GetRepo("dev-tool").Url, 18, nil)
	assert.Nil(t, err)
	util.PrintJson(mr)

}

func TestNewGitRepo(t *testing.T) {
	config := GetConfig().GetGitLabConfig("https://git.internal.yunify.com/chenyu/doc/")
	util.PrintJson(config)

	git, err := gitlab.NewClient(config.Token, gitlab.WithBaseURL(config.BaseUrl))
	assert.Nil(t, err)

	pp, _, err := git.Projects.GetProject("chenyu/doc", nil)
	assert.Nil(t, err)

	util.PrintJson(pp)

}

func TestGitAutoJiraID(t *testing.T) {
	//dir := "/Users/yu/code/git.internal.yunify.com/chenyu/doc"
	dir := "/Users/yu/code/git.internal.yunify.com/bi/pitrix-wh-daemon"
	jiraID, cType, cm := AutoJiraID(dir, nil, "")
	fmt.Println(jiraID, cType, cm)

	jiraID, cType, cm = AutoJiraID(dir, nil, "a67f4bae")
	fmt.Println(jiraID, cType, cm)
}

func TestAutoBranch(t *testing.T) {
	branch := AutoBranch("/Users/yu/code/git.internal.yunify.com/bi/pitrix-wh-daemon")
	fmt.Println(branch)
}

func TestMatch(t *testing.T) {
	input := "feature/PRODUCT-1234:test"

	// 定义正则表达式模式，匹配 PRODUCT- 后跟一串数字或字母
	re := regexp.MustCompile(`PRODUCT-\w+`)

	// 查找匹配项
	match := re.FindString(input)

	// 输出匹配项
	fmt.Println("Matched:", match)

}

func TestMatch1(t *testing.T) {
	re := regexp.MustCompile(`^[a-z0-9]+$`)
	fmt.Println(re.MatchString("aaa"))
}

func TestGetMergeRequest(t *testing.T) {
	config := GetConfig().GetGitLabConfig("https://git.internal.yunify.com/chenyu/doc/")
	util.PrintJson(config)

	git, err := gitlab.NewClient(config.Token, gitlab.WithBaseURL(config.BaseUrl))
	assert.Nil(t, err)

	res, _, err := git.MergeRequests.ListMergeRequests(&gitlab.ListMergeRequestsOptions{
		State:        stringPtr("opened"),
		SourceBranch: stringPtr("VM-1888_x_staging"),
		TargetBranch: stringPtr("staging"),
	})
	assert.Nil(t, err)

	fmt.Println(res[0])

}

func TestGetMergeRequest1(t *testing.T) {
	config := GetConfig().GetGitLabConfig("https://git.internal.yunify.com/Simon/pitrix-billing")
	util.PrintJson(config)

	git, err := gitlab.NewClient(config.Token, gitlab.WithBaseURL(config.BaseUrl))
	assert.Nil(t, err)
	res, _, err := git.MergeRequests.GetMergeRequest("Simon/pitrix-billing", 1123, nil)
	assert.Nil(t, err)
	fmt.Println(res)

}
