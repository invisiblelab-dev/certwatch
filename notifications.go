package certwatch

type Email struct {
	Username string `yml:"username"`
	Password string `yml:"password"`
	SMTPHost string `yml:"smtphost"`
	Port     int    `yml:"port"`
	From     string `yml:"from"`
	To       string `yml:"to"`
}
