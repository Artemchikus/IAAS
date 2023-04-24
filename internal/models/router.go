package models

import "time"

type Router struct {
	CreatedAt           time.Time            `json:"created_at"`
	Description         string               `json:"description"`
	FlavorID            string               `json:"flavor_id"`
	ID                  string               `json:"id"`
	Interfaces          []string             `json:"interfaces_info"`
	ExternalGatewayInfo *ExternalGatewayInfo `json:"external_gateway_info"`
	Name                string               `json:"name"`
	ProjectID           string               `json:"project_id"`
	Routes              []string             `json:"routes"`
	Status              string               `json:"status"`
	UpdatedAt           time.Time            `json:"updated_at"`
}

type ExternalGatewayInfo struct {
	NetworkID string `json:"network_id"`
}

func NewRouter(Description, Name, ExternalNetworkId string) *Router {
	info := &ExternalGatewayInfo{
		NetworkID: ExternalNetworkId,
	}

	return &Router{
		Description:         Description,
		Name:                Name,
		ExternalGatewayInfo: info,
	}
}

// {
// 	"admin_state_up": true,
// 	"availability_zone_hints": [],
// 	"availability_zones": [],
// 	"created_at": "2023-04-03T14:35:08Z",
// 	"description": "",
// 	"enable_ndp_proxy": null,
// 	"external_gateway_info": {
// 	  "network_id": "de4171b0-67c2-4140-81b6-066fb97d7f06",
// 	  "external_fixed_ips": [
// 		{
// 		  "subnet_id": "deba4b1e-2388-41a0-bca7-8ab95f2cbc59",
// 		  "ip_address": "192.168.122.237"
// 		}
// 	  ],
// 	  "enable_snat": true
// 	},
// 	"flavor_id": null,
// 	"id": "2efe00cc-f740-4b39-8364-b417340652b8",
// 	"interfaces_info": [
// 	  {
// 		"port_id": "bcdf7ade-2c55-4d8d-aeb8-714d2cb4b336",
// 		"ip_address": "192.168.100.1",
// 		"subnet_id": "94a7c47e-517e-452b-8abc-9b9958663d10"
// 	  }
// 	],
// 	"name": "router01",
// 	"project_id": "09522dddf25648d5bc30307cf3bf5f72",
// 	"revision_number": 4,
// 	"routes": [],
// 	"status": "ACTIVE",
// 	"tags": [],
// 	"tenant_id": "09522dddf25648d5bc30307cf3bf5f72",
// 	"updated_at": "2023-04-03T14:37:55Z"
//   }
