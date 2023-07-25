package certwatch

type ConfigFile struct {
	Domains       []Domain `yml:"domains"`
	Refresh       int      `yml:"refresh"`
	Notifications struct {
		Email `yml:"email"`
	} `yml:"notifications"`
}
