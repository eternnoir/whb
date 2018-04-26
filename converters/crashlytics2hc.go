package converters

import "github.com/labstack/echo"

type Crashlytics2HC struct {
}

func (ch *Crashlytics2HC) SourceName() string {
	return "crashlytics"
}

func (ch *Crashlytics2HC) TargetName() string {
	return "hangoutchat"
}
func (ch *Crashlytics2HC) Process(c echo.Context) error {
	return nil
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
