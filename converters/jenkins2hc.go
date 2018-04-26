package converters

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func NewJenkins2HC() *Jenkins2HC {
	tmpl, err := template.New("j2hct").Parse(Jenkins2HCTemplate)
	if err != nil {
		panic(err)
	}
	return &Jenkins2HC{template: tmpl}
}

type Jenkins2HC struct {
	template *template.Template
}

func (jh *Jenkins2HC) SourceName() string {
	return "jenkins"
}

func (jh *Jenkins2HC) TargetName() string {
	return "hangoutchat"
}

func (jh *Jenkins2HC) Process(c echo.Context) error {
	p := new(JenkinsPayload)
	if err := c.Bind(p); err != nil {
		logrus.WithError(err).Error("Decode payload fail.")
		return err
	}
	params := ParseHangoutChatParams(c)
	logrus.WithField("params", params).WithField("req", p).Info("Receive")
	msg, err := jh.toMsg(p)
	if err != nil {
		return err
	}
	if err := SendHangoutText(params, msg); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, p)
}

func (jh *Jenkins2HC) toMsg(jp *JenkinsPayload) (string, error) {
	var tpl bytes.Buffer
	if err := jh.template.Execute(&tpl, jp); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

type JenkinsPayload struct {
	Build struct {
		Artifacts struct {
			Asgard_standalone_jar struct {
				Archive string `json:"archive"`
				S3      string `json:"s3"`
			} `json:"asgard-standalone.jar"`
			Asgard_war struct {
				Archive string `json:"archive"`
			} `json:"asgard.war"`
		} `json:"artifacts"`
		FullURL string `json:"full_url"`
		Number  int    `json:"number"`
		Phase   string `json:"phase"`
		Scm     struct {
			Branch string `json:"branch"`
			Commit string `json:"commit"`
			URL    string `json:"url"`
		} `json:"scm"`
		Status string `json:"status"`
		URL    string `json:"url"`
	} `json:"build"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (jp *JenkinsPayload) IsFail() bool {
	return jp.Build.Status == ""
}

var Jenkins2HCTemplate = `
{
          {{if eq .Build.Status "FAILURE"}}
  "text": "<users/all> Jenkins Job Fail!!!!!",
          {{end}}
  "cards": [
    {
      "header": {
        "title": "{{.Name}}",
        "subtitle": "{{.Build.Status}}",
        {{if eq .Build.Status "SUCCESS"}}
         "imageUrl": "http://ci.gotyourpoint.com:8080/static/813b537a/images/headshot.png"
        {{else if eq .Build.Status ""}}
         "imageUrl": "http://ci.gotyourpoint.com:8080/static/813b537a/images/headshot.png"
        {{else}}
          "imageUrl": "https://i2.wp.com/halclan.net/wp-content/uploads/2017/08/jenkins-logo.jpg"
        {{end}}
      },
      "sections": [
        {
          "widgets": [
              {
                "keyValue": {
                  "topLabel": "Build Number",
                  "content": "{{.Build.Number}}"
                  }
              },
              {
                "keyValue": {
                  "topLabel": "Phase",
                  "content": "{{.Build.Phase}}"
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
                            "url": "{{.Build.FullURL}}"
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
