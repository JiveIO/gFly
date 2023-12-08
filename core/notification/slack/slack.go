package slack

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func New(webHook string) Slack {
	return Slack{
		WebHook: webHook,
	}
}

type Slack struct {
	WebHook string
}

func (slack Slack) Notify(text string, username ...string) error {
	json := buildRequest(text, username...)

	req, err := http.NewRequest("POST", slack.WebHook, strings.NewReader(json))
	if err != nil {
		return fmt.Errorf("can't connect to host %s: %s", slack.WebHook, err.Error())
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	return err
}

func buildRequest(text string, username ...string) string {
	msg := fmt.Sprintf(`{"text":" %s "}`, text)

	if len(username) > 0 {
		msg = fmt.Sprintf(`{"text":" <@%s> %s "}`, username[0], text)
	}

	return msg
}
