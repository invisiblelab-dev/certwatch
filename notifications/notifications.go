package notifications

import (
	"fmt"
	"math"

	"github.com/invisiblelab-dev/certwatch"
)

type MessageData struct {
	Messages []string
}

type Notifier interface {
	Notify(title string, message MessageData, recipients ...string) error
	Recipient() string
}

type NotifierService struct {
	notifiers []Notifier
}

func NewNotifierService() *NotifierService {
	return &NotifierService{}
}

func (n *NotifierService) Append(notifier Notifier) {
	n.notifiers = append(n.notifiers, notifier)
}

func (n *NotifierService) Notify(title string, message MessageData) error {
	for _, notifier := range n.notifiers {
		if err := notifier.Notify(title, message, notifier.Recipient()); err != nil {
			return fmt.Errorf("failed to notify: %w", err)
		}
	}

	return nil
}

// TODO: review and move strings into embed template
func ComposeMessage(domainDeadlines []certwatch.DomainDeadline) (MessageData, error) {
	var message MessageData
	for _, domainDeadline := range domainDeadlines {
		if domainDeadline.OnDeadline {
			domain := domainDeadline.Domain
			if domainDeadline.DaysTillDeadline <= 0 {
				days := int64(math.Abs(domainDeadline.DaysTillDeadline))
				message.Messages = append(message.Messages, fmt.Sprintf("- %s certificate has expired %d days ago.", domain, days))
			} else {
				days := int64(domainDeadline.DaysTillDeadline)
				message.Messages = append(message.Messages, fmt.Sprintf("- %s certificate expires in %d days.", domain, days))
			}
		}
	}

	return message, nil
}
