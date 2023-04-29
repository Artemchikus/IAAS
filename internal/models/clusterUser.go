package models

type ClusterUser struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DomainID    string `json:"domain_id" toml:"domain_id"`
	ProjectID   string `json:"project_id" toml:"project_id"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Description string `json:"description"`
	ClusterID   int    `json:"cluster_id"`
	AccountID   int    `json:"account_id"`
	CluserRole  string `json:"cluser_role"`
}

func NewClusterUser(name, email, password, projectId, domainId, description string) *ClusterUser {
	return &ClusterUser{
		Name:        name,
		Email:       email,
		Password:    password,
		ProjectID:   projectId,
		DomainID:    domainId,
		Description: description,
	}
}
