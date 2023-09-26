package cmd

import (
	"bee-flow/cmd/rss"
	"os"

	"github.com/spf13/cobra"
)

var longDesc = `
______     ______     ______     ______   __         ______     __     __
/\  == \   /\  ___\   /\  ___\   /\  ___\ /\ \       /\  __ \   /\ \  _ \ \
\ \  __<   \ \  __\   \ \  __\   \ \  __\ \ \ \____  \ \ \/\ \  \ \ \/ ".\ \
 \ \_____\  \ \_____\  \ \_____\  \ \_\    \ \_____\  \ \_____\  \ \__/".~\_\
  \/_____/   \/_____/   \/_____/   \/_/     \/_____/   \/_____/   \/_/   \/_/


源码地址： https://github.com/xifo-wu/bee-flow
文档地址： https://blog.xifo.in
`

var (
	CfgFile string
	rootCmd = &cobra.Command{
		Use:   "beeflow",
		Short: "一个简陋的 qbittorrent 扩展命令行",
		Long:  longDesc,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	cobra.OnInitialize(InitConfig)

	rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "配置文件 (默认地址为 $HOME/.config/bee-flow/config.yaml)")

	rootCmd.AddCommand(rss.RssCmd)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
