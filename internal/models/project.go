package models

type Project struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Enabled     bool     `json:"enabled"`
	DomainID    string   `json:"domain_id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Options     *Options `json:"options"`
}

type Options struct{}

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
