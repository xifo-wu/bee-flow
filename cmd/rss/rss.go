/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package rss

import (
	"github.com/spf13/cobra"
)

// rssCmd represents the rss command
var RssCmd = &cobra.Command{
	Use:   "rss",
	Short: "订阅相关功能",
}

func init() {
	RssCmd.AddCommand(rssAddCmd)
}
