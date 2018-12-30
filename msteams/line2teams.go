package msteams

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/labstack/echo"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/sirupsen/logrus"
)

func NewLine2Teams() *Line2Teams {
	tmpl, err := template.New("line2hc").Parse(LINE2HCMsgTemplae)
	if err != nil {
		panic(err)
	}
	return &Line2Teams{
		template: tmpl,
	}
}

type Line2Teams struct {
	template *template.Template
}

func (lh *Line2Teams) SourceName() string {
	return "line"
}

func (lh *Line2Teams) TargetName() string {
	return "teams"
}

func (lh *Line2Teams) Process(c echo.Context) error {
	secret, token := getLineToken(c)
	bot, err := linebot.New(secret, token)
	if err != nil {
		logrus.WithError(err).Error("Init line bot fail.")
		return err
	}
	params := ParseTeamsParams(c)
	logrus.WithField("params", params).Info("Line Converter receive")

	events, err := bot.ParseRequest(c.Request())
	if err != nil {
		logrus.WithError(err).Error("Parse Line event error")
		if err == linebot.ErrInvalidSignature {
			return c.NoContent(400)
		}
		return c.NoContent(500)
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			sourceName, err := lh.getEventSourceName(bot, event)
			if err != nil {
				logrus.WithError(err).Error("Line bot get source DisplayName fail")
			}
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				return lh.porcessMessage(c, &LineMsg{
					DisplayName: sourceName,
					Event:       event,
					TextMsg:     message,
				}, params)
			case *linebot.StickerMessage:
				return lh.porcessMessage(c, &LineMsg{
					DisplayName: sourceName,
					Event:       event,
					StickerMsg:  message,
				}, params)
			}
		}
	}
	return nil
}

func (lh *Line2Teams) getEventSourceName(bot *linebot.Client, event linebot.Event) (string, error) {
	user, err := bot.GetProfile(event.Source.UserID).Do()
	if err != nil {
		return "", err
	}
	return user.DisplayName, nil
}

func (lh *Line2Teams) porcessMessage(c echo.Context, msg *LineMsg, params *TeamsParams) error {
	var tpl bytes.Buffer
	if err := lh.template.Execute(&tpl, msg); err != nil {
		return err
	}
	if err := SendTeamsText(params, tpl.String()); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// TODO refactor to line package
type LineMsg struct {
	DisplayName string
	Event       linebot.Event
	TextMsg     *linebot.TextMessage
	StickerMsg  *linebot.StickerMessage
}

func getLineToken(c echo.Context) (secret, token string) {
	secret = c.QueryParam("l_secret")
	token = c.QueryParam("l_token")
	return
}

var LINE2HCMsgTemplae = `
{
    "@context": "https://schema.org/extensions",
    "@type": "MessageCard",
    "themeColor": "0072C6",
    "title": "{{.DisplayName}}",
    "text": "{{.TextMsg.Text}}"
}
`
