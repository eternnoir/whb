package converters

import (
	"fmt"
	"net/http"

	"github.com/eternnoir/whb/conmgr"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func init() {
	conmgr.DefuaultConverterMgr.RegConverter(&Jenkins2HC{})
}

var Jenkins2HCTemplate = "[%s] %s - #%d - %s (<%s|Open>)"

type Jenkins2HC struct {
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
	params := &HangoutChatParams{
		Spaces: c.QueryParam("gh_spaces"),
		Key:    c.QueryParam("gh_key"),
		Token:  c.QueryParam("gh_token"),
	}

	logrus.WithField("params", params).WithField("req", p).Info("Receive")
	msg := jh.toMsg(p)
	if err := SendHangoutTextMesssage(params, msg); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, p)
}

func (jh *Jenkins2HC) toMsg(jp *JenkinsPayload) string {
	return fmt.Sprintf(Jenkins2HCTemplate,
		jp.Build.Phase,
		jp.Name,
		jp.Build.Number,
		jp.Build.Status,
		jp.Build.FullURL,
	)
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
