package pkg

import (
	"bee-flow/pkg/logger"
	"bee-flow/pkg/torrent"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/spf13/viper"
)

func Rename(torrentName string, path string) []string {
	// 读取存储的数据
	item := viper.GetStringMap(fmt.Sprintf("data.%s", path))
	isItemOk := len(item) > 0

	// qb 下载保持到的位置，可能是文件也可能是文件夹
	downloadPath := filepath.Join(path, torrentName)
	var files []string
	file, err := os.Stat(downloadPath)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	if !file.Mode().IsRegular() {
		logger.Logger.Info("下载目标为目录，开始移动到父目录")
		files = append(files, MoveToParentDir(path, torrentName)...)
	} else {
		files = append(files, torrentName)
	}

	newNames := make([]string, 0, len(files))

	for _, file := range files {
		logger.Logger.Info("重命名开始：" + file)
		logger.Logger.Info(fmt.Sprintf("是否存在于数据库中：%t", isItemOk))
		if isItemOk {
			newName := RenameFile(item, path, file)
			oldPath := filepath.Join(path, file)
			newPath := filepath.Join(path, newName)

			os.Rename(oldPath, newPath)
			newNames = append(newNames, newPath)

			logger.Logger.Info("新文件名：" + newPath)
		} else {
			newPath := filepath.Join(path, file)
			newNames = append(newNames, newPath)
			logger.Logger.Info("新文件名：" + newPath)
		}
	}

	return newNames
}

func MoveToParentDir(path string, torrentName string) []string {
	dirPath := filepath.Join(path, torrentName)
	// 获取目录下的所有文件
	files, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	filesNames := make([]string, 0, len(files))

	// 遍历所有文件
	for _, file := range files {
		// 判断是否为文件
		if !file.IsDir() {
			// 拼接旧文件路径和新文件路径
			oldFilePath := filepath.Join(dirPath, file.Name())

			newFilePath := filepath.Join(filepath.Dir(dirPath), file.Name())

			// 移动文件到上级目录
			err := os.Rename(oldFilePath, newFilePath)
			if err != nil {
				panic(err)
			}

			filesNames = append(filesNames, file.Name())
		}
	}

	return filesNames
}

func RenameFile(item map[string]interface{}, path string, filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return filename
	}

	// zip 文件跳过重命名
	collectionRe := regexp.MustCompile(`(?i)(\d+-\d+|第\d+-\d+集|合集)`)
	if collectionRe.FindString(filename) != "" {
		// 暂时什么都不处理
		return filename
	}

	assRe := regexp.MustCompile(`(?i)\.(zip|7z)$`)
	if assRe.MatchString(filename) {
		return filename
	}

	mode, ok := item["mode"].(float64)
	logger.Logger.Info(fmt.Sprintf("重命名模式：%f", mode))
	if !ok || mode == 0 {
		return filename
	}

	log.Println("filename: ", filename)
	if mode == 1 {
		return RenameMode1(path, filename, item)
	}

	if mode == 2 {
		return RenameMode2(path, filename, item)
	}

	return filename
}

func RenameMode1(path string, filename string, item map[string]interface{}) string {
	offset := item["offset"]
	SE := GenerateSeasonAndEpisode(path, filename, offset.(float64))
	ext := filepath.Ext(filename)

	return fmt.Sprintf("%s %s%s", item["name"], SE, ext)
}

func RenameMode2(path string, filename string, item map[string]interface{}) string {
	offset := item["offset"]
	SE := GenerateSeasonAndEpisode(path, filename, offset.(float64))
	ext := filepath.Ext(filename)
	name := item["name"]
	multiVersion := GenerateMultiVersion(item)

	if ext == ".ass" {
		languageCode := ""
		// 字幕重命名
		// 寻找语言关键字 CHS CHT
		// 正则表达式模式
		re := regexp.MustCompile(`(?i)\.(CHS|CHT|SC|TC)\.ass`)
		match := re.FindStringSubmatch(filename)
		if len(match) > 1 {
			languageCode = match[1]
			fmt.Printf("语言代码: %s\n", languageCode)
		}

		if languageCode != "" {
			return fmt.Sprintf("%s - %s - %s.%s.%s", name, SE, multiVersion, languageCode, ext)
		}
	}

	return fmt.Sprintf("%s - %s - %s%s", name, SE, multiVersion, ext)
}

func GenerateSeasonAndEpisode(path string, filename string, offset float64) string {
	log.Print(path, "pathpathpath")
	standardTitleRe := regexp.MustCompile(`S\d+E\d+`)
	// 查找第一个匹配的子字符串
	match := standardTitleRe.FindString(filename)

	if match != "" {
		return match
	}

	re := regexp.MustCompile(`(?i)(?:Season (\d+)|Season(\d+)|S(\d+))`) // 编译正则表达式
	// 匹配季数
	matches := re.FindStringSubmatch(path)
	var seasonNumber int
	// 如果有匹配结果，打印第一个匹配分组中的数字部分
	if len(matches) > 0 {
		// 第一个匹配分组中的数字部分
		seasonNumber, _ = strconv.Atoi(matches[1])
	}

	torrentInfo := torrent.TorrentNameParse(filename, path)

	if math.Floor(torrentInfo.Episode) == torrentInfo.Episode {
		return fmt.Sprintf("S%02dE%02d", seasonNumber, int(torrentInfo.Episode)-int(offset))
	}

	return fmt.Sprintf("S%02dE%04.1f", seasonNumber, torrentInfo.Episode-offset)
}

func GenerateMultiVersion(item map[string]interface{}) string {
	multiVersion, ok := (item["multi_version"]).(string)

	if !ok {
		group, ok := item["group"]
		if ok {
			multiVersion += fmt.Sprintf("%s", group)
		}

		resolution, ok := item["resolution"]
		if ok {
			multiVersion += fmt.Sprintf(".%s", resolution)
		}

		subtitle, ok := item["subtitle"]
		if ok {
			multiVersion += fmt.Sprintf(".%s", subtitle)
		}
	}

	return multiVersion
}
