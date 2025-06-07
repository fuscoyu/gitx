package cmd

import (
	"fmt"
	"github.com/goeoeo/gitx/repo"
	"github.com/spf13/cobra"
)

var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "info",
	Run: func(cmd *cobra.Command, args []string) {
		config := repo.GetConfig(configPath).Init()
		r, err := config.CurrentRepo()
		config.CheckErr(err)
		fmt.Println("repo url is: ", r.Url)
	},
}
