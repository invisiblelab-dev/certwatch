package certwatch

type ConfigFile struct {
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
