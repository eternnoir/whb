package msteams

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"fmt"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

var (
	TeamsWebhook = "https://outlook.office.com/webhook/%s/IncomingWebhook/%s"
)

type TeamsParams struct {
	Webhook          string
	IncommingWebhook string
}

func ParseTeamsParams(c echo.Context) *TeamsParams {
	params := &TeamsParams{
		Webhook:          c.QueryParam("teams_wh"),
		IncommingWebhook: c.QueryParam("teams_in_wh"),
	}
	return params
}

func SendTeamsText(params *TeamsParams, data string) error {
	url := fmt.Sprintf(TeamsWebhook, params.Webhook, params.IncommingWebhook)
	logrus.WithField("url", url).WithField("params", params).WithField("msg", data).Info("Send teams webhook request")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
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
