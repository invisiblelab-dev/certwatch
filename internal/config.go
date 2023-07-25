package certwatch

type ConfigFile struct {
	Domains       []Domain `yml:"domains"`
	Notifications struct {
		Email Email `yml:"email"`
	} `yml:"notifications"`
}
