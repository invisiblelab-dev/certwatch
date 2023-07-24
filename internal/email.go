package certwatch

type Email struct {
	Mailtrap struct {
		Username string `yml:"username"`
		Password string `yml:"password"`
		SmtpHost string `yml:"smtpHost"`
	} `yml:"mailtrap"`
	From string `yml:"from"`
	To   string `yml:"to"`
}
