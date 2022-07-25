/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ccc469/sshctl/internal"
	"github.com/ccc469/sshctl/internal/login"

	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "SSH登录远程主机",
	Long: fmt.Sprintf(`SSH登录远程主机，支持密码、密钥文件登录


Your config dir %s
`, internal.HomePath),
	Run: func(cmd *cobra.Command, args []string) {
		login.Run()
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&internal.Ip, "ip", "i", "127.0.0.1", "服务器地址")
	loginCmd.Flags().IntVarP(&internal.Port, "port", "P", 22, "SSH端口")
	loginCmd.Flags().StringVarP(&internal.User, "user", "u", "root", "用户名")
	loginCmd.Flags().BoolVar(&internal.Pass, "pass", true, "验证方式（username/password or ssh key")
	loginCmd.Flags().BoolVar(&internal.Save, "save", true, "是否保存连接信息")
	loginCmd.Flags().StringVarP(&internal.PrivateKey, "private-key", "k", strings.Join([]string{os.Getenv("HOME"), ".ssh", "id_rsa"}, internal.Delimiter), "private key path.")
	loginCmd.Flags().StringVarP(&internal.AliasName, "alias-name", "n", "", "连接名, 默认为IP")

	loginCmd.MarkFlagRequired("ip")
}
