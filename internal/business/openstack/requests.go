package openstack

type GetTokenRequest struct {
	Auth *Identity `json:"auth"`
}

type Identity struct {
	Methods  []string  `json:"methods"`
	Password *Password `json:"password"`
	Token    *Token    `json:"token"`
}

type Password struct {
	User *User `json:"user"`
}

type User struct {
	Name     string  `json:"name"`
	Domain   *Domain `json:"domain"`
	Password string  `json:"password"`
}

type Domain struct {
	ID string `json:"id"`
}

type Token struct {
	ID string `json:"id"`
}
