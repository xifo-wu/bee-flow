package torrent

import (
	"fmt"
	"regexp"
	"strconv"
)

type TorrentInfo struct {
	Episode float64
}

func TorrentNameParse(torrentName string, path string) TorrentInfo {
	var torrentInfo TorrentInfo
	const (
		Pattern1 = `(?i)(.*) - ((\d{1,4}\.\d{1,2})|(\d{1,4}))(?:v\d{1,2})?(?: )?(?:END)?(.*)`
		Pattern2 = `(?i)(.*)[\[\ E?P](\d{1,4}|\d{1,4}\.\d{1,2})(?:v\d{1,2})?(?: )?(?:END)?[\]\ ](.*)`
		Pattern3 = `(?i)(.*)\[(?:第)?(\d*\.*\d*)[话集話](?:END)?\](.*)`
		Pattern4 = `(?i)(.*)第?(\d*\.*\d*)[话話集](?:END)?(.*)`
		Pattern5 = `(?i)(.*)EP(\d{1,4})(?:v\d{1,2})?(?: )?(?:END)?(.*)`
	)

	// 遍历正则表达式模式
	patterns := []string{Pattern1, Pattern2, Pattern3, Pattern4, Pattern5}
	for _, pattern := range patterns {
		// 编译正则表达式
		re := regexp.MustCompile(pattern)

		// 查找匹配的子字符串及其子匹配项
		matches := re.FindStringSubmatch(torrentName)

		// 如果有匹配结果，打印第一个匹配分组中的数字部分
		if len(matches) > 0 {
			// 第一个匹配分组中的数字部分
			torrentInfo.Episode, _ = strconv.ParseFloat(matches[2], 64) // 使用索引 2 获取第一个匹配数字
			fmt.Printf("匹配模式: %s\n", pattern)
			fmt.Printf("匹配的数字部分: %f\n", torrentInfo.Episode)
			break // 返回第一个匹配的数字后退出循环
		}
	}

	return torrentInfo
}
