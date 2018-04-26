package converters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

var (
	HangoutChatWebHookUrl = "https://chat.googleapis.com/v1/spaces/%s/messages?key=%s&token=%s"
)

type HangoutChatParams struct {
	Spaces string `query:"gh_spaces"`
	Key    string `query:"gh_key"`
	Token  string `query:"gh_token"`
}

type HangoutChatMessage struct {
	Text string `json:"text"`
}

func SendHangoutTextMesssage(params *HangoutChatParams, message string) error {
	url := fmt.Sprintf(HangoutChatWebHookUrl, params.Spaces, params.Key, params.Token)
	logrus.WithField("url", url).WithField("params", params).WithField("msg", message).Info("Send hangout chant webhook request")
	jsonStr, _ := json.Marshal(HangoutChatMessage{Text: message})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	logrus.Info("response Body:", string(body))
	return nil
}