package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化 BeeFlow 程序",
	Long:  `使用 Init 命令将会初始化 BeeFlow 程序，它将自动生成配置文件`,
	Run: func(cmd *cobra.Command, args []string) {
		InitRun(args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func InitRun(args []string) {
	InitConfig()
}

func InitConfig() {
	d := color.New(color.FgRed, color.Bold)

	// 不存在配置文件时将会创建
	home, err := homedir.Dir()
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	filePath := filepath.Join(home, ".config", "bee-flow", "config.yaml")
	currentDir := filepath.Dir(filePath)

	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(currentDir)
	}

	if !checkFileExists(filePath) {
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			d.Printf("创建目录时发生错误: %s", err)
			return
		}

		viper.SetDefault("QB_ADDRESS", "http://127.0.0.1:8080")
		viper.SetDefault("QB_USERNAME", "admin")
		viper.SetDefault("QB_PASSWORD", "")
		viper.SetDefault("TMDB_V4_TOKEN", "")
		viper.SetDefault("TELEGRAM_BOT_TOKEN", "")
		viper.SetDefault("TELEGRAM_CHAT_ID", "")
		viper.SetDefault("TELEGRAM_CHANNEL_ID", "")
		viper.SetDefault("HDHIVE_USERNAME", "")
		viper.SetDefault("HDHIVE_PASSWORD", "")
		viper.SetDefault("HDHIVE_TOKEN", "")
		viper.SetDefault("BACKUP_PATH", "")

		e := viper.WriteConfigAs(filePath)
		fmt.Println(currentDir, "currentDir", e, "Eee")

		g := color.New(color.FgGreen, color.Bold)
		g.Println("配置初始化成功")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)

	return !errors.Is(error, os.ErrNotExist)
}
