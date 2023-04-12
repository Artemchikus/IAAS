package models

import "time"

type Network struct {
	AvailabilityZones []string  `json:"availability_zones"`
	CreatedAt         time.Time `json:"created_at"`
	Description       string    `json:"description"`
	ID                string    `json:"id"`
	MTU               int       `json:"mtu"`
	Name              string    `json:"name"`
	ProjectID         string    `json:"project_id"`
	Status            string    `json:"status"`
	Subnets           []string  `json:"subnets"`
	Tags              []string  `json:"tags"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// {
// 	"admin_state_up": true,
// 	"availability_zone_hints": [],
// 	"availability_zones": [],
// 	"created_at": "2023-04-03T14:35:24Z",
// 	"description": "",
// 	"dns_domain": null,
// 	"id": "439e51bd-f03f-44c5-9006-72a68d401e4a",
// 	"ipv4_address_scope": null,
// 	"ipv6_address_scope": null,
// 	"is_default": null,
// 	"is_vlan_transparent": null,
// 	"mtu": 1442,
// 	"name": "private",
// 	"port_security_enabled": true,
// 	"project_id": "09522dddf25648d5bc30307cf3bf5f72",
// 	"provider:network_type": "geneve",
// 	"provider:physical_network": null,
// 	"provider:segmentation_id": 58474,
// 	"qos_policy_id": null,
// 	"revision_number": 3,
// 	"router:external": false,
// 	"segments": null,
// 	"shared": true,
// 	"status": "ACTIVE",
// 	"subnets": [
// 	  "94a7c47e-517e-452b-8abc-9b9958663d10"
// 	],
// 	"tags": [],
// 	"updated_at": "2023-04-03T14:49:52Z"
//   }
