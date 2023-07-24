package certwatch

type ConfigFile struct {
	Domains       []Domain `yml:"domains"`
	Roots         bool     `yml:"roots"`
	Notifications struct {
		Email Email `yml:"email"`
	} `yml:"notifications"`
}
