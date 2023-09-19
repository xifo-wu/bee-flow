package cmd

import (
	"bee-flow/pkg"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	torrentName string
	category    string
	rootPath    string
	contentPath string
	savePath    string
	torrentSize string
	infoHash    string
)

// moveCmd represents the move command
var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "使用 Rclone 移动文件",
	Long:  `使用 Rclone 移动文件，达到文件上传的目的`,
	Run: func(cmd *cobra.Command, args []string) {
		MoveCmdFunc(args)
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)

	rootCmd.PersistentFlags().StringVarP(&torrentName, "name", "n", "", "Torrent 名称")
	rootCmd.PersistentFlags().StringVarP(&category, "category", "l", "", "分类")
	rootCmd.PersistentFlags().StringVarP(&contentPath, "contentPath", "f", "", "内容路径（与多文件 torrent 的根目录相同）")
	rootCmd.PersistentFlags().StringVarP(&rootPath, "rootPath", "r", "", "根目录（第一个 torrent 的子目录路径）")
	rootCmd.PersistentFlags().StringVarP(&savePath, "savePath", "d", "", "保存路径")
	rootCmd.PersistentFlags().StringVarP(&torrentSize, "torrentSize", "z", "", "Torrent 大小（字节）")
	rootCmd.PersistentFlags().StringVarP(&infoHash, "infoHash", "i", "", "T信息哈希值 v1")
}

func MoveCmdFunc(args []string) {
	// 根据分类执行不同命令

	// isRss := strings.Contains(category, "RSS")
	// isTV := strings.Contains(category, "TV")
	// isMovie := strings.Contains(category, "Movie")

	if !strings.Contains(category, "BeeFlow") {
		red := color.New(color.FgRed)
		red.Print("非 BeeFlow 相关分类不运行")

		return
	}

	// 如果只是 BeeFlow 就只执行备份
	if category == "BeeFlow" {
		moveFile(savePath, viper.GetString("backup_path"))
		return
	}
}

func moveFile(src string, dst string) {
	originalPath := filepath.Join(src, torrentName)
	targetPath := filepath.Join(dst, torrentName)

	cmd := exec.Command("rclone", "moveto", "-v", "-P", originalPath, targetPath)

	if err := cmd.Run(); err != nil {
		log.Println("Rclone Error")
		log.Fatal(err)
	}

	pkg.Notification(torrentName + "下载完成")
}
