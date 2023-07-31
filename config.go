package certwatch

import (
	"fmt"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Slack struct {
	Enabled bool `yml:"enabled"`
	Webhook string
}

type Email struct {
	Enabled  bool   `yml:"enabled"`
	Username string `yml:"username"`
	Password string `yml:"password"`
	SMTPHost string `yml:"smtphost"`
	Port     int    `yml:"port"`
	From     string `yml:"from"`
	To       string `yml:"to"`
}

type Stdout struct {
	Enabled bool `yml:"enabled"`
}

type Domain struct {
	Name      string `yaml:"name"`
	Threshold int    `yaml:"threshold"`
}

type Cache struct {
	Path    string `yml:"path"`
	Refresh int    `yml:"refresh"`
}

type Config struct {
	Domains       []Domain `yml:"domains"`
	Cache         `yml:"cache"`
	Notifications struct {
		Email  `yml:"email"`
		Slack  `yml:"slack"`
		Stdout `yml:"stdout"`
	} `yml:"notifications"`
}

func (c *Config) UnmarshalYAML(payload []byte) error {
	re := regexp.MustCompile(`^(?:[_a-z0-9](?:[_a-z0-9-]{0,61}[a-z0-9])?\.)+(?:[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?)?$`)

	if err := yaml.Unmarshal(payload, &c); err != nil {
		return fmt.Errorf("error parsing yaml: %w", err)
	}

	for _, domain := range c.Domains {
		if !re.Match([]byte(domain.Name)) {
			return fmt.Errorf("invalid domain name [%s]", domain.Name)
		}
	}

	return nil
}
