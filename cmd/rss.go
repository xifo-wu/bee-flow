/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// rssCmd represents the rss command
var rssCmd = &cobra.Command{
	Use:   "rss",
	Short: "订阅功能",
}

func init() {
	rootCmd.AddCommand(rssCmd)
}
