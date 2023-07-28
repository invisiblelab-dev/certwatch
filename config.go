package certwatch

type Config struct {
	Domains       []Domain `yml:"domains"`
	Refresh       int      `yml:"refresh"`
	Notifications struct {
		Email `yml:"email"`
		Slack `yml:"slack"`
	} `yml:"notifications"`
}

type Slack struct {
	Webhook string
}

type Email struct {
	Username string `yml:"username"`
	Password string `yml:"password"`
	SMTPHost string `yml:"smtphost"`
	Port     int    `yml:"port"`
	From     string `yml:"from"`
	To       string `yml:"to"`
}
