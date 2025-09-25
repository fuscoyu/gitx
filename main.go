package main

import (
	"fmt"
	"os"

	"github.com/goeoeo/gitx/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitx",
	Short: "GitX - Enhanced Git workflow tool",
	Long:  `GitX is a command-line tool that enhances Git workflow with JIRA integration and automated processes.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VersionInfo())
	},
}

func main() {
	rootCmd.AddCommand(cmd.PushCmd, cmd.PullCmd, cmd.JiraCmd, cmd.InitCmd, cmd.InfoCmd, cmd.HookCmd, versionCmd)
	
	if err := rootCmd.Execute(); err != nil {
		logrus.Debugf("run cmd err:%s", err)
		os.Exit(1)
	}
}
