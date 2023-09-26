package pkg

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/spf13/viper"
)

func FindOrCreateData() map[string]map[string]interface{} {
	myData := make(map[string]map[string]interface{})
	dataPath := viper.GetString("data_path")
	file, err := os.OpenFile(dataPath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		log.Println("无法打开/创建文件:", err)
		return myData
	}

	defer file.Close()

	// 尝试解码JSON数据，如果文件为空或无效，则初始化一个新的结构体
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&myData); err != nil {
		if err != io.EOF {
			log.Println("解码JSON失败:", err)
			return myData
		}
	}

	return myData
}

func UpdateData(content map[string]map[string]interface{}) bool {
	dataPath := viper.GetString("data_path")

	// 将内容编码为 JSON 格式
	jsonData, err := json.Marshal(content)
	if err != nil {
		log.Println("编码JSON失败:", err)
		return false
	}

	// 使用 os.WriteFile 将内容写入文件，替换整个文件
	if err := os.WriteFile(dataPath, jsonData, os.ModePerm); err != nil {
		log.Println("写入文件失败:", err)
		return false
	}

	return true
}
