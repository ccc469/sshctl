/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/ccc469/sshctl/internal"
	"github.com/ccc469/sshctl/internal/server"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "服务连接",
	Long:  `管理服务器连接信息`,
	Run: func(cmd *cobra.Command, args []string) {
		if res := server.Run(); !res {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().BoolVar(&internal.ClearServer, "clear", false, "清空服务器连接信息")
	serverCmd.Flags().StringVar(&internal.RemoveAliasName, "del", "", "连接名")
	serverCmd.Flags().BoolVar(&internal.ShowServer, "show", false, "查看所有连接信息")

}
