package notifications

import (
	"fmt"
	"math"
	"strings"

	"github.com/invisiblelab-dev/certwatch"
)

type Notifier interface {
	Notify(title string, message string, recipients ...string) error
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

func (n *NotifierService) Notify(title string, message string) error {
	for _, notifier := range n.notifiers {
		if err := notifier.Notify(title, message, notifier.Recipient()); err != nil {
			return fmt.Errorf("failed to notify: %w", err)
		}
	}

	return nil
}

// TODO: review and move strings into embed template
func ComposeMessage(domainDeadlines []certwatch.DomainDeadline) (string, error) {
	message := strings.Builder{}
	for _, domainDeadline := range domainDeadlines {
		if domainDeadline.OnDeadline {
			domain, err := certwatch.RemoveHTTPS(domainDeadline.Domain)
			if err != nil {
				return "", fmt.Errorf("failed to remove HTTPS form URL: %w", err)
			}
			if domainDeadline.DaysTillDeadline <= 0 {
				days := int64(math.Abs(domainDeadline.DaysTillDeadline))
				message.WriteString(fmt.Sprintf("- %s certificate has expired %d days ago.", domain, days))
			} else {
				days := int64(domainDeadline.DaysTillDeadline)
				message.WriteString(fmt.Sprintf("- %s certificate expires in %d days.", domain, days))
			}
			message.WriteString("\n")
		}
	}

	return message.String(), nil
}
