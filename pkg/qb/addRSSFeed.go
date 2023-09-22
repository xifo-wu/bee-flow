package qb

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/spf13/viper"
)

func AddRSSFeed(feedUrl string, path string) bool {
	url := fmt.Sprintf("%s/api/v2/rss/addFeed", viper.GetString("QB_ADDRESS"))
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("url", feedUrl)
	if path != "" {
		_ = writer.WriteField("path", path)
	}

	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Add("Cookie", "SID="+viper.GetString("QB_SID"))

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fmt.Println(string(body), res.StatusCode)
	return res.StatusCode > 200
}
