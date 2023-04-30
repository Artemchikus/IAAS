package models

type Role struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

func NewRole(Description, Name string) *Role {
	return &Role{
		Description: Description,
		Name:        Name,
	}
}

// {
// 	"description": null,
// 	"domain_id": null,
// 	"id": "891f0037561a47bea621fac592225737",
// 	"name": "reader",
// 	"options": {
// 	  "immutable": true
// 	}
//   }
