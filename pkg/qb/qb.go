package qb

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/spf13/viper"
)

func Login() {
	url := fmt.Sprintf("%s/api/v2/auth/login", viper.GetString("QB_ADDRESS"))
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("username", viper.GetString("QB_USERNAME"))
	_ = writer.WriteField("password", viper.GetString("QB_PASSWORD"))
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	sidValue := ""
	// 检查响应中的Set-Cookie头字段
	cookies := res.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "SID" {
			sidValue = cookie.Value
			fmt.Println("SID:", sidValue)
			break
		}
	}

	viper.Set("QB_SID", sidValue)
	viper.WriteConfig()
	fmt.Println("QB Login:", string(body))
}

func CreateCategory(category string, savePath string) {
	url := fmt.Sprintf("%s/api/v2/torrents/createCategory", viper.GetString("QB_ADDRESS"))
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("category", category)
	if savePath != "" {
		_ = writer.WriteField("savePath", savePath)
	}

	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "SID="+viper.GetString("QB_SID"))

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}
