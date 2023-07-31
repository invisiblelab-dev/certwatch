package factory

import (
	"github.com/invisiblelab-dev/certwatch"
	"github.com/invisiblelab-dev/certwatch/notifications"
)

type Factory struct {
	NotifierService func() notifications.NotifierService
	Config          certwatch.Config
}

func NewNotifierService(f *Factory) func() notifications.NotifierService {
	return func() notifications.NotifierService {
		service := notifications.NewNotifierService()

		if f.Config.Notifications.Stdout.Enabled {
			service.Append(notifications.NewStdoutNotifier())
		}

		if f.Config.Notifications.Email.Enabled {
			cfg := notifications.EmailNotifierConfig{
				Login:    f.Config.Notifications.Email.Username,
				Password: f.Config.Notifications.Email.Password,
				Host:     f.Config.Notifications.Email.SMTPHost,
				Port:     f.Config.Notifications.Email.Port,
				From:     f.Config.Notifications.Email.From,
				To:       f.Config.Notifications.Email.To,
			}
			service.Append(notifications.NewEmailNotifier(cfg))
		}

		if f.Config.Notifications.Slack.Enabled {
			service.Append(notifications.NewSlackNotifier(f.Config.Notifications.Slack.Webhook))
		}

		return *service
	}
}
