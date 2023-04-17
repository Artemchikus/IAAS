package models

type Subnet struct {
	AllocationPools []*AllocationPool `json:"allocation_pools"`
	CIDR            string            `json:"cidr"`
	CreatedAt       string            `json:"created_at"`
	Description     string            `json:"description"`
	EnableDHCP      bool              `json:"enable_dhcp"`
	GatewayIp       string            `json:"gateway_ip"`
	ID              string            `json:"id"`
	IpVersion       int               `json:"ip_version"`
	Name            string            `json:"name"`
	NetworkID       string            `json:"network_id"`
	ProjectID       string            `json:"project_id"`
	Tags            []string          `json:"tags"`
	UpdatedAt       string            `json:"updated_at"`
}

type AllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// {
// 	"allocation_pools": [
// 	  {
// 		"start": "192.168.100.2",
// 		"end": "192.168.100.254"
// 	  }
// 	],
// 	"cidr": "192.168.100.0/24",
// 	"created_at": "2023-04-03T14:35:50Z",
// 	"description": "",
// 	"dns_nameservers": [],
// 	"dns_publish_fixed_ip": null,
// 	"enable_dhcp": true,
// 	"gateway_ip": "192.168.100.1",
// 	"host_routes": [],
// 	"id": "94a7c47e-517e-452b-8abc-9b9958663d10",
// 	"ip_version": 4,
// 	"ipv6_address_mode": null,
// 	"ipv6_ra_mode": null,
// 	"name": "private-subnet",
// 	"network_id": "439e51bd-f03f-44c5-9006-72a68d401e4a",
// 	"project_id": "09522dddf25648d5bc30307cf3bf5f72",
// 	"revision_number": 0,
// 	"segment_id": null,
// 	"service_types": [],
// 	"subnetpool_id": null,
// 	"tags": [],
// 	"updated_at": "2023-04-03T14:35:50Z"
//   }
