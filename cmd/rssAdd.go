package cmd

import (
	"bee-flow/pkg"
	"bee-flow/pkg/qb"
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	path        string
	mode        int
	tmdb        string
	rssSavePath string
)

var rssAddCmd = &cobra.Command{
	Use:   "add",
	Short: "[暂不支持相同路径]添加一个新的 RSS 订阅",
	Run: func(cmd *cobra.Command, args []string) {
		AddRSSRun(args)
	},
}

func init() {
	rssCmd.AddCommand(rssAddCmd)

	rssAddCmd.Flags().IntVarP(&mode, "mode", "m", 0, "重命名模式")
	rssAddCmd.Flags().StringVarP(&path, "path", "p", "", "RSS订阅文件夹 (e.g. ‘The Pirate Bay\\Top100\\Video’)")
	rssAddCmd.Flags().StringVarP(&tmdb, "tmdb", "t", "", "TMDB 的链接地址")
	rssAddCmd.Flags().StringVarP(&rssSavePath, "savePath", "s", "", "下载保存的地址，会自动加上配置内的路径")
}

// AddRSSRun 向 qbittorrent 添加一个订阅，并创建自动下载器等
func AddRSSRun(args []string) {
	argsLen := len(args)
	if argsLen == 0 {
		log.Println("订阅地址必填")
		return
	}

	if mode <= 4 && tmdb == "" {
		log.Println("重命名模式为 1 到 4 时 TMDB 链接必须填写。 请添加 -t 标识")
		return
	}

	url := args[0]
	savePath := filepath.Join(viper.GetString("backup_path"), rssSavePath)

	data := map[string]interface{}{
		"enabled":          true,
		"mustNotContain":   "合集",
		"useRegex":         true,
		"affectedFeeds":    [1]string{url},
		"assignedCategory": "BeeFlow",
		"savePath":         savePath,
		"addPaused":        true,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	qb.Login()
	isError := qb.AddRSSFeed(url, path)
	if isError {
		return
	}

	db := pkg.FindOrCreateData()

	// 获取对应的名字和年份

	db[savePath] = map[string]interface{}{
		"mode": mode,
	}

	qb.SetRSSRule(url, string(jsonData))
	// 存储一份到 db json
	pkg.UpdateData(db)
}
