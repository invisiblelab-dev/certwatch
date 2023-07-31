package certwatch

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
