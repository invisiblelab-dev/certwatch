package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Text struct {
	Type string   `json:"type"`
	Text TextBody `json:"text"`
}

type TextBody struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji,omitempty"`
}

type SlackNotifier struct {
	Webhook string
}

func NewSlackNotifier(webhook string) *SlackNotifier {
	return &SlackNotifier{webhook}
}

func (s *SlackNotifier) Notify(title string, message MessageData, recipients ...string) error {
	msg := strings.Join(message.Messages, "\n")
	payload, err := json.Marshal(s.blocks(title, msg))
	if err != nil {
		return fmt.Errorf("failed to marshal slack payload: %w", err)
	}

	if err := s.post(recipients[0], payload); err != nil {
		return fmt.Errorf("failed to send slack message: %w", err)
	}

	return nil
}

func (s *SlackNotifier) Recipient() string {
	return s.Webhook
}

func (s *SlackNotifier) blocks(title string, message string) map[string]any {
	return map[string]any{
		"blocks": []Text{
			{
				Type: "section",
				Text: TextBody{
					Type: "mrkdwn",
					Text: fmt.Sprintf("🚨 *%s*", title),
				},
			},
			{
				Type: "section",
				Text: TextBody{
					Type: "mrkdwn",
					Text: message,
				},
			},
		},
	}
}

func (s *SlackNotifier) post(recipient string, payload []byte) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, recipient, bytes.NewBuffer(payload))
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
