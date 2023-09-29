package version

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var VERSION string

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看程序当前版本",
	Long:  `查看程序当前版本`,
	Run: func(cmd *cobra.Command, args []string) {
		notice := color.New(color.BgHiBlack, color.FgHiWhite).PrintlnFunc()
		notice("Don't forget this...")

		c := color.New(color.FgCyan, color.Bold)
		c.Println(VERSION)
	},
}

func init() {
}
