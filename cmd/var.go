package cmd

var (
	branchList           string //目标分支逗号分隔
	planTgtBranchList    string //计划要推的分支列表,逗号分隔
	action               string //jira 执行的动作
	configPath           string //配置文件
	project              string //项目地址
	jiraID               string //jiraID
	force                bool   //忽略本地记录，cherry-pick所有commit
	debug                bool   //开启debug日志
	disableAutoMergeHook bool   // 自动合并后是否执行hook
)
