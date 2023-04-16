package openstack

type GetTokenRequest struct {
	Auth *Auth `json:"auth"`
}

type Auth struct {
	Identity *Identity `json:"identity"`
	Scope    *Scope    `json:"scope"`
}
type Identity struct {
	Methods  []string  `json:"methods"`
	Password *Password `json:"password"`
}

type Scope struct {
	Project *GetTokenProject `json:"project"`
}

type GetTokenProject struct {
	Name   string  `json:"name"`
	Domain *Domain `json:"domain"`
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

type CreateProjectRequest struct {
	Project *CreateProject `json:"project"`
}

type CreateProject struct {
	DomainID    string `json:"domain_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}
