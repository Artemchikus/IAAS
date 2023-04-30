package models

import "time"

type SecurityGroup struct {
	CreatedAt   time.Time       `json:"created_at"`
	Description string          `json:"description"`
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	ProjectID   string          `json:"project_id"`
	Rules       []*SecurityRule `json:"rules"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func NewSecurityGroup(Description, Name string) *SecurityGroup {
	return &SecurityGroup{
		Description: Description,
		Name:        Name,
	}
}

// {
// 	"created_at": "2023-04-03T13:24:55Z",
// 	"description": "secgroup01",
// 	"id": "1f2f9c20-92bd-4e27-b6b3-c8530cd0d49a",
// 	"name": "secgroup01",
// 	"project_id": "09522dddf25648d5bc30307cf3bf5f72",
// 	"revision_number": 4,
// 	"rules": [
// 	  {
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
// 	  {
// 		"id": "242a07d6-e3ea-49ec-93b8-6ba634bd96d8",
// 		"tenant_id": "09522dddf25648d5bc30307cf3bf5f72",
// 		"security_group_id": "1f2f9c20-92bd-4e27-b6b3-c8530cd0d49a",
// 		"ethertype": "IPv4",
// 		"direction": "ingress",
// 		"protocol": "tcp",
// 		"port_range_min": 22,
// 		"port_range_max": 22,
// 		"remote_ip_prefix": "0.0.0.0/0",
// 		"remote_address_group_id": null,
// 		"normalized_cidr": "0.0.0.0/0",
// 		"remote_group_id": null,
// 		"standard_attr_id": 54,
// 		"description": "",
// 		"tags": [],
// 		"created_at": "2023-04-03T15:10:15Z",
// 		"updated_at": "2023-04-03T15:10:15Z",
// 		"revision_number": 0,
// 		"project_id": "09522dddf25648d5bc30307cf3bf5f72"
// 	  },
// 	  {
// 		"id": "5d89fd73-6f79-4f60-9409-64f93f518deb",
// 		"tenant_id": "09522dddf25648d5bc30307cf3bf5f72",
// 		"security_group_id": "1f2f9c20-92bd-4e27-b6b3-c8530cd0d49a",
// 		"ethertype": "IPv4",
// 		"direction": "egress",
// 		"protocol": null,
// 		"port_range_min": null,
// 		"port_range_max": null,
// 		"remote_ip_prefix": null,
// 		"remote_address_group_id": null,
// 		"normalized_cidr": null,
// 		"remote_group_id": null,
// 		"standard_attr_id": 18,
// 		"description": null,
// 		"tags": [],
// 		"created_at": "2023-04-03T13:24:55Z",
// 		"updated_at": "2023-04-03T13:24:55Z",
// 		"revision_number": 0,
// 		"project_id": "09522dddf25648d5bc30307cf3bf5f72"
// 	  },
// 	  {
// 		"id": "7dd42ea9-8102-484a-9a50-4afec1986719",
// 		"tenant_id": "09522dddf25648d5bc30307cf3bf5f72",
// 		"security_group_id": "1f2f9c20-92bd-4e27-b6b3-c8530cd0d49a",
// 		"ethertype": "IPv4",
// 		"direction": "ingress",
// 		"protocol": "icmp",
// 		"port_range_min": null,
// 		"port_range_max": null,
// 		"remote_ip_prefix": "0.0.0.0/0",
// 		"remote_address_group_id": null,
// 		"normalized_cidr": "0.0.0.0/0",
// 		"remote_group_id": null,
// 		"standard_attr_id": 53,
// 		"description": "",
// 		"tags": [],
// 		"created_at": "2023-04-03T15:10:03Z",
// 		"updated_at": "2023-04-03T15:10:03Z",
// 		"revision_number": 0,
// 		"project_id": "09522dddf25648d5bc30307cf3bf5f72"
// 	  },
// 	  {
// 		"id": "f3434ca5-f757-4ba4-9d60-13a0d7e96a02",
// 		"tenant_id": "09522dddf25648d5bc30307cf3bf5f72",
// 		"security_group_id": "1f2f9c20-92bd-4e27-b6b3-c8530cd0d49a",
// 		"ethertype": "IPv4",
// 		"direction": "ingress",
// 		"protocol": "tcp",
// 		"port_range_min": 80,
// 		"port_range_max": 80,
// 		"remote_ip_prefix": "0.0.0.0/0",
// 		"remote_address_group_id": null,
// 		"normalized_cidr": "0.0.0.0/0",
// 		"remote_group_id": null,
// 		"standard_attr_id": 81,
// 		"description": "",
// 		"tags": [],
// 		"created_at": "2023-04-03T16:48:13Z",
// 		"updated_at": "2023-04-03T16:48:13Z",
// 		"revision_number": 0,
// 		"project_id": "09522dddf25648d5bc30307cf3bf5f72"
// 	  }
// 	],
// 	"stateful": true,
// 	"tags": [],
// 	"updated_at": "2023-04-03T16:48:13Z"
//   }
