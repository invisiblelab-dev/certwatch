package certwatch

type Email struct {
	Username string `yml:"username"`
	Password string `yml:"password"`
	SmtpHost string `yml:"smtphost"`
	From     string `yml:"from"`
	To       string `yml:"to"`
}
