package certwatch

type ConfigFile struct {
	Domains       []Domain `yml:"domains"`
	Notifications struct {
		Email `yml:"email"`
	} `yml:"notifications"`
}
