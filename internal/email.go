package certwatch

type Email struct {
	Mailtrap struct {
		Username string `yml:"username"`
		Password string `yml:"password"`
		SmtpHost string `yml:"smtphost"`
	} `yml:"mailtrap"`
	From string `yml:"from"`
	To   string `yml:"to"`
}
