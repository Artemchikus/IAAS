package models

type SecurityRule struct {
	ID              string   `json:"id"`
	ProjectId       string   `json:"project_id"`
	SecurityGroupID string   `json:"security_group_id"`
	Ethertype       string   `json:"ethertype"`
	Direction       string   `json:"direction"`
	Protocol        string   `json:"protocol"`
	PortRangeMax    int      `json:"port_range_max"`
	PortRangeMin    int      `json:"port_range_min"`
	RemoteIpPrefix  string   `json:"remote_ip_prefix"`
	Description     string   `json:"description"`
	Tags            []string `json:"tags"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

// {
// 		"id": "075d629e-702a-4928-97ba-f37a30b1dc58",
// 		"tenant_id": "09522dddf25648d5bc30307cf3bf5f72",
// 		"security_group_id": "1f2f9c20-92bd-4e27-b6b3-c8530cd0d49a",
// 		"ethertype": "IPv6",
// 		"direction": "egress",
// 		"protocol": null,
// 		"port_range_min": null,
// 		"port_range_max": null,
// 		"remote_ip_prefix": null,
// 		"remote_address_group_id": null,
// 		"normalized_cidr": null,
// 		"remote_group_id": null,
// 		"standard_attr_id": 19,
// 		"description": null,
// 		"tags": [],
// 		"created_at": "2023-04-03T13:24:55Z",
// 		"updated_at": "2023-04-03T13:24:55Z",
// 		"revision_number": 0,
// 		"project_id": "09522dddf25648d5bc30307cf3bf5f72"
// 	  },
