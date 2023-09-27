package cmd

import (
	"bee-flow/pkg"
	"bee-flow/pkg/hdhive"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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

	moveCmd.Flags().StringVarP(&torrentName, "name", "n", "", "Torrent 名称")
	moveCmd.Flags().StringVarP(&category, "category", "l", "", "分类")
	moveCmd.Flags().StringVarP(&contentPath, "contentPath", "f", "", "内容路径（与多文件 torrent 的根目录相同）")
	moveCmd.Flags().StringVarP(&rootPath, "rootPath", "r", "", "根目录（第一个 torrent 的子目录路径）")
	moveCmd.Flags().StringVarP(&savePath, "savePath", "d", "", "保存路径")
	moveCmd.Flags().StringVarP(&torrentSize, "torrentSize", "z", "", "Torrent 大小（字节）")
	moveCmd.Flags().StringVarP(&infoHash, "infoHash", "i", "", "T信息哈希值 v1")
}

func MoveCmdFunc(args []string) {
	if !strings.Contains(category, "BeeFlow") {
		red := color.New(color.FgRed)
		red.Print("非 BeeFlow 相关分类不运行")

		return
	}

	db := pkg.FindOrCreateData()
	item, isItemOk := db[savePath]

	newNames := pkg.Rename(torrentName, savePath)

	for _, file := range newNames {
		backPath := viper.GetString("backup_path")

		if isItemOk {
			ok := false
			backPath, ok = item["backupPath"].(string)
			if !ok {
				backPath = strings.Replace(file, viper.GetString("save_base_path"), viper.GetString("backup_path"), 1)
			}
		}

		// 使用 filepath.Base 获取文件名称（包括扩展名）
		fullFileName := filepath.Base(file)
		backPath = filepath.Join(backPath, fullFileName)

		moveFile(file, backPath)

		telegramChannelID := viper.GetString("telegram_channel_id")
		if telegramChannelID != "" {
			sleepDuration := 1 * time.Minute
			time.Sleep(sleepDuration)

			NotificationTGChannel(file, item)
		}
	}
}

func NotificationTGChannel(src string, data map[string]interface{}) {
	// 使用 filepath.Base 获取文件名称（包括扩展名）
	fullFileName := filepath.Base(src)

	videoRe := regexp.MustCompile(`\.(mp4|mov|avi|wmv|mkv|flv|webm|vob|rmvb|mpg|mpeg)$`)
	if !videoRe.MatchString(strings.ToLower(fullFileName)) {
		return
	}

	hdhive.Notification(src, data)

	log.Println("同步完成")
}

func moveFile(src string, dst string) {
	cmd := exec.Command("rclone", "moveto", "-v", "-P", src, dst)

	if err := cmd.Run(); err != nil {
		log.Println("Rclone Error")
		log.Fatal(err)
	}

	pkg.Notification(src + "下载完成")
}
