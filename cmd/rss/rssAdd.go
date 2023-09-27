package rss

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
	backupPath     string
	renameMode     int
	name           string
	year           int
	season         int
	rssSavePath    string
	group          string
	offset         int
	resolution     string
	subtitle       string
	multiVersion   string
	hdhiveShareId  string
	mustNotContain string
)

var rssAddCmd = &cobra.Command{
	Use:   "add",
	Short: "[暂不支持相同路径]添加一个新的 RSS 订阅",
	Long: `[暂不支持相同路径]添加一个新的 RSS 订阅.
更多订阅规则前往 QB 控制台进行修改
	`,
	Run: func(cmd *cobra.Command, args []string) {
		AddRSSRun(args)
	},
}

func init() {
	rssAddCmd.Flags().IntVarP(&renameMode, "mode", "m", 0, "重命名模式")
	rssAddCmd.Flags().StringVarP(&name, "name", "n", "", "订阅影视名称")
	rssAddCmd.Flags().IntVarP(&year, "year", "y", 0, "订阅影视年份")
	rssAddCmd.Flags().IntVarP(&season, "season", "s", 0, "订阅影视季度")
	rssAddCmd.Flags().StringVarP(&resolution, "resolution", "", "", "分辨率，部分重命名格式用得到")
	rssAddCmd.Flags().StringVarP(&group, "group", "", "", "字幕组，部分重命名格式用得到")
	rssAddCmd.Flags().StringVarP(&backupPath, "path", "p", "", "备份路径, 不填为全局配置")
	rssAddCmd.Flags().StringVarP(&subtitle, "subtitle", "", "", "字幕组，部分重命名格式用得到")
	rssAddCmd.Flags().IntVarP(&offset, "offset", "o", 0, "有些命名会从多集开始，添加这个参数后会自动减去集数")
	rssAddCmd.Flags().StringVarP(&rssSavePath, "savePath", "r", "", "下载保存的地址，会自动加上配置内的路径")
	rssAddCmd.Flags().StringVarP(&multiVersion, "multiVersion", "", "", "多版本信息，mode2 专用")
	rssAddCmd.Flags().StringVarP(&hdhiveShareId, "hdhiveShareId", "", "", "通过影巢通知到频道参数，影巢分享记录ID")
	rssAddCmd.Flags().StringVarP(&mustNotContain, "mustNotContain", "", "", "RSS 订阅不可包含（填写正则表达式）")

}

// AddRSSRun 向 qbittorrent 添加一个订阅，并创建自动下载器等
func AddRSSRun(args []string) {
	argsLen := len(args)
	if argsLen == 0 {
		log.Println("订阅地址必填")
		return
	}

	if renameMode <= 4 && renameMode > 0 && (name == "" || year == 0) {
		log.Println("重命名模式为 1 到 4 时 TMDB 链接必须填写。 请添加 -name 和 -year 标识")
		return
	}

	url := args[0]
	savePath := filepath.Join(viper.GetString("save_base_path"), rssSavePath)

	if mustNotContain == "" {
		mustNotContain = "合集"
	}

	data := map[string]interface{}{
		"enabled":          true,
		"mustNotContain":   mustNotContain,
		"useRegex":         true,
		"affectedFeeds":    [1]string{url},
		"assignedCategory": "BeeFlow",
		"savePath":         savePath,
		"addPaused":        false,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	qb.Login()

	log.Println(backupPath, "backupPath")
	isError := qb.AddRSSFeed(url, backupPath)

	if isError {
		return
	}

	db := pkg.FindOrCreateData()

	// 获取对应的名字和年份

	db[savePath] = map[string]interface{}{
		"mode":          renameMode,
		"name":          name,
		"year":          year,
		"season":        season,
		"rssSavePath":   rssSavePath,
		"backupPath":    backupPath,
		"group":         group,
		"offset":        offset,
		"resolution":    resolution,
		"subtitle":      subtitle,
		"multiVersion":  multiVersion,
		"hdhiveShareId": hdhiveShareId,
	}

	qb.SetRSSRule(url, string(jsonData))
	// 存储一份到 db json
	pkg.UpdateData(db)
}
