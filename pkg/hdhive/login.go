package hdhive

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

func Login() string {
	url := "https://www.hdhive.org/api/v1/login"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{"username":"%s","password":"%s"}`, viper.GetString("hdhive_username"), viper.GetString("hdhive_password")))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		// JSON 解析失败
		return ""
	}

	log.Println(response, "login response")

	success, ok := response["success"]
	if !ok || !success.(bool) {
		return ""
	}

	meta, ok := (response["meta"]).(map[string]interface{})
	if !ok {
		return ""
	}

	return meta["access_token"].(string)
}
