package openstack

type GetTokenRequest struct {
	Auth *GetTokenAuth `json:"auth"`
}
type RefreshTokenRequest struct {
	Auth *RefreshTokenAuth `json:"auth"`
}
type RefreshTokenAuth struct {
	Identity *TokenIdentity `json:"identity"`
}
type TokenIdentity struct {
	Methods []string `json:"methods"`
	Token   *Token   `json:"password"`
}
type GetTokenAuth struct {
	Identity *PasswordIdentity `json:"identity"`
	Scope    *Scope            `json:"scope"`
}
type PasswordIdentity struct {
	Methods  []string  `json:"methods"`
	Password *Password `json:"password"`
}
type Scope struct {
	Project *GetTokenProject `json:"project"`
}
type GetTokenProject struct {
	ID     string  `json:"id"`
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
type CreateUserRequest struct {
	User *CreateUser `json:"user"`
}
type CreateUser struct {
	Name      string `json:"name"`
	DomainID  string `json:"domain_id"`
	ProjectID string `json:"default_project_id"`
	Password  string `json:"password"`
	Enabled   bool   `json:"enabled"`
	Email     string `json:"email"`
}
type CretaeImageRequest struct {
	DiskFormat      string `json:"disk_format"`
	ContainerFormat string `json:"container_format"`
	Name            string `json:"name"`
	Visibility      string `json:"visibility"`
}

type CreateFlavorRequest struct {
	Flavor *Flavor `json:"flavor"`
}

type Flavor struct {
	VCPUs      int     `json:"vcpus"`
	Disk       int     `json:"disk"`
	Name       string  `json:"name"`
	RAM        int     `json:"ram"`
	Ephemeral  int     `json:"OS-FLV-EXT-DATA:ephemeral"`
	IsPublic   bool    `json:"os-flavor-access:is_public"`
	Swap       string  `json:"swap"`
	RXTXFactor float32 `json:"rxtx_factor"`
}
