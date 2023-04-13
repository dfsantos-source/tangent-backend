package tangent

type Util struct {
	token string
}

func CreateUtil(token string) *Util {
	return &Util{
		token: token,
	}
}

func (u Util) GetToken() string {
	return u.token
}
