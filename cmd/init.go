package cmd

import (
	"fmt"
	"github.com/goeoeo/gitx/repo"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

// configContentTpl 配置文件模板
var configContentTpl = `
log_level: 0  #debug:5

patch:
  plan_tgt_branch_list: [ "dev","qa" ] #默认计划要推的分支
  branch_alias:  #分支别名
    v6.0: QCE_V6.0-20220630
    v6.1: QCE_V6.1-20221230
    v6.2: QCE_V6.2-20231230

gitLab_configs:
  - base_url: https://git.yunify.com
    token: "gitlab Access Tokens 用于自动创建mr,合并mr" 
  - base_url: https://git.internal.yunify.com
    token: "gitlab Access Tokens 用于自动创建mr,合并mr"

repo:
  doc:  #git项目简称
    url: "https://git.internal.yunify.com/chenyu/doc" #项目https地址
    path: "" #项目本地根路径,需要使用绝对路径
    create_mr: true #自动创建mr
    auto_merge_branch_list: ["dev","qa"] #自动合并的分支


`
var try bool
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化项目配置文件",
	Run: func(cmd *cobra.Command, args []string) {
		if try {
			config := repo.GetConfig()
			config.Print()

			return
		}
		err := createPatchConfig()
		checkErr(err)

		return
	},
}

func init() {
	InitCmd.Flags().BoolVarP(&try, "try", "t", false, "打印配置文件")

}
func createPatchConfig() (err error) {
	var (
		homeDir string
	)

	if homeDir, err = os.UserHomeDir(); err != nil {
		return
	}

	dir := homeDir + "/.patch"
	if _, err = os.Stat(dir); err != nil {
		if _, err = repo.ExecCmd(homeDir, "mkdir", ".patch"); err != nil {
			return
		}
	}

	configPath := homeDir + "/.patch/config.yaml"

	// 文件不存在，创建
	if _, err = os.Stat(configPath); err == nil {
		fmt.Println("已存在配置：", configPath)
		return
	}

	fmt.Println("请配置你的信息，路径为：", configPath)

	err = ioutil.WriteFile(configPath, []byte(configContentTpl), os.ModePerm)
	return
}
