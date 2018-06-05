package hangoutschat

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func NewCrashlytics2HC() *Crashlytics2HC {
	tmpl, err := template.New("c2hct").Parse(Crashlytics2HCTemplate)
	if err != nil {
		panic(err)
	}
	return &Crashlytics2HC{template: tmpl}
}

type Crashlytics2HC struct {
	template *template.Template
}

func (ch *Crashlytics2HC) SourceName() string {
	return "crashlytics"
}

func (ch *Crashlytics2HC) TargetName() string {
	return "hangoutschat"
}
func (ch *Crashlytics2HC) Process(c echo.Context) error {
	p := new(CrashlyticsPayload)
	if err := c.Bind(p); err != nil {
		logrus.WithError(err).Error("Decode payload fail.")
		return err
	}
	params := ParseHangoutChatParams(c)
	logrus.WithField("params", params).WithField("req", p).Info("Receive")
	if p.Event == "verification" {
		return c.NoContent(http.StatusOK)
	}
	var tpl bytes.Buffer
	if err := ch.template.Execute(&tpl, p); err != nil {
		return err
	}
	if err := SendHangoutText(params, tpl.String()); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

type CrashlyticsPayload struct {
	Event   string `json:"event"`
	Payload struct {
		CrashesCount         int    `json:"crashes_count"`
		DisplayID            int    `json:"display_id"`
		ImpactLevel          int    `json:"impact_level"`
		ImpactedDevicesCount int    `json:"impacted_devices_count"`
		Method               string `json:"method"`
		Title                string `json:"title"`
		URL                  string `json:"url"`
	} `json:"payload"`
	PayloadType string `json:"payload_type"`
}

var Crashlytics2HCTemplate = `
{
  "text": "<users/all> Crashlytics!!!!!",
  "cards": [
    {
      "header": {
        "title": "{{.Payload.Title}}",
        "subtitle": "{{.Payload.Method}}",
		"imageUrl": "https://pbs.twimg.com/profile_images/438136190760284160/Q3LJPRGx_400x400.png"
      },
      "sections": [
        {
          "widgets": [
            {
              "keyValue": {
                "topLabel": "CrashesCount",
                "content": "{{.Payload.CrashesCount}}"
              }
            },
            {
              "keyValue": {
                "topLabel": "ImpactLevel",
                "content": "{{.Payload.ImpactLevel}}"
              }
            },
            {
              "keyValue": {
                "topLabel": "ImpactedDevicesCount",
                "content": "{{.Payload.ImpactedDevicesCount}}"
              }
            },
            {
              "keyValue": {
                "topLabel": "DisplayID",
                "content": "{{.Payload.DisplayID}}"
              }
            }
          ]
        },
        {
          "widgets": [
            {
              "buttons": [
                {
                  "textButton": {
                    "text": "OPEN URL",
                    "onClick": {
                      "openLink": {
                        "url": "{{.Payload.URL}}"
                      }
                    }
                  }
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
`
