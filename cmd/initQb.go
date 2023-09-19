/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bee-flow/pkg/qb"

	"github.com/spf13/cobra"
)

// initQbCmd represents the initQb command
var initQbCmd = &cobra.Command{
	Use:   "initQb",
	Short: "初始化 QB 相关资料，需完成 QB 配置",
	Run: func(cmd *cobra.Command, args []string) {
		qb.Login()
		qb.CreateCategory("BeeFlow", "")
		qb.CreateCategory("BeeFlow+TV", "")
		qb.CreateCategory("BeeFlow+RSS+TV", "")
		qb.CreateCategory("BeeFlow+RSS+Movie", "")
	},
}

func init() {
	rootCmd.AddCommand(initQbCmd)
}
