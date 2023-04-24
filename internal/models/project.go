package models

type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DomainID    string `json:"domain_id"`
	Description string `json:"description"`
}

type Options struct{}

func NewProject(name, description string) *Project {
	return &Project{
		Name:        name,
		Description: description,
		DomainID:    "default",
	}
}

// {
// 	"description": "Service Project",
// 	"domain_id": "default",
// 	"enabled": true,
// 	"id": "2e83a8e5362247f79c8d86980ab2a216",
// 	"is_domain": false,
// 	"name": "service",
// 	"options": {},
// 	"parent_id": "default",
// 	"tags": []
//   }
