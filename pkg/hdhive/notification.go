package hdhive

import (
	"bee-flow/pkg"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/spf13/viper"
)

const WEBAPI = "https://www.hdhive.org/api/v1/manager/telegram/auto-message"

var (
	client = &http.Client{}
)

func Notification(src string, data map[string]interface{}) {
	telegramBotToken := viper.GetString("telegram_bot_token")
	telegramChannelID := viper.GetString("telegram_channel_id")

	resourceId, ok := data["hdhive_share_id"]
	if !ok {
		return
	}

	// 匹配 SxxExx
	standardTitleRe := regexp.MustCompile(`S\d+E\d+`)
	// 符合 S01E01 时直接返回文件名，不需要重命名
	SE := standardTitleRe.FindString(src)

	multiVersion := pkg.GenerateMultiVersion(data)

	remarkFormat := SE + " - " + multiVersion

	payload := make(map[string]interface{})
	payload["chat_id"] = telegramChannelID
	payload["telegram_bot_token"] = telegramBotToken
	payload["remark"] = remarkFormat
	payload["resource_id"] = resourceId
	payload["share_size"] = 0

	log.Println(payload, "payload")

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
		return
	}

	req, err := http.NewRequest("POST", WEBAPI, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatal(err)
	}

	// 偷懒每次都登录
	// TODO 使用缓存 Token
	apiToken := Login()
	if apiToken == "" {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		// JSON 解析失败
		return
	}

	fmt.Println(response, "Res")
}
