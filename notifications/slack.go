package notifications

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	certwatch "github.com/invisiblelab-dev/certwatch"
)

func SendSlack(subject string, slackHook certwatch.Slack) error {
	url := "https://hooks.slack.com/services/" + slackHook.Webhook

	var jsonStr = []byte(fmt.Sprintf(`{"text": "%s"}`, subject))
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return fmt.Errorf("failed to setup slack request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post slack message: %w", err)
	}

	defer resp.Body.Close()

	return nil
}
