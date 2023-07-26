package notifications

import (
	"bytes"
	"fmt"
	"net/http"

	certwatch "github.com/invisiblelab-dev/certwatch"
)

func SendSlack(subject string, slackHook certwatch.Slack) error {

	url := "https://hooks.slack.com/services/" + slackHook.Webhook

	var jsonStr = []byte(fmt.Sprintf(`{"text": "%s"}`, subject))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)
	return err
}
