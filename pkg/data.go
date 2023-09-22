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
	file, err := os.OpenFile(dataPath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		log.Println("无法打开/创建文件:", err)
		return false
	}

	file.Seek(0, 0) // 将文件指针移回文件开头
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(content); err != nil {
		log.Println("编码JSON失败:", err)
		return false
	}

	return true
}
