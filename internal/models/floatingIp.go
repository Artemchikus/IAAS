package models

import "time"

type FloatingIp struct {
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	FixedIp     string    `json:"fixed_ip_address"`
	FolatingIp  string    `json:"floating_ip_address"`
	ID          string    `json:"id"`
	NetworkID   string    `json:"floating_network_id"`
	PortID      string    `json:"port_id"`
	ProjectID   string    `json:"project_id"`
	RouterID    string    `json:"router_id"`
	Status      string    `json:"status"`
	SubnetID    string    `json:"subnet_id"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewFloatingIP(NetworkID, description string) *FloatingIp {
	return &FloatingIp{
		NetworkID:   NetworkID,
		Description: description,
	}
}

// {
// 	"created_at": "2023-04-03T15:09:34Z",
// 	"description": "",
// 	"dns_domain": "",
// 	"dns_name": "",
// 	"fixed_ip_address": "192.168.100.31",
// 	"floating_ip_address": "192.168.122.213",
// 	"floating_network_id": "de4171b0-67c2-4140-81b6-066fb97d7f06",
// 	"id": "23d844f5-c6b8-4f15-8621-4894b110ada0",
// 	"name": "192.168.122.213",
// 	"port_details": "admin_state_up='True', device_id='1a639b0f-30fa-47c7-9ba3-9512ce1583d6', device_owner='compute:nova', mac_address='fa:16:3e:83:0d:c1', name='', network_id='439e51bd-f03f-44c5-9006-72a68d401e4a', status='DOWN'",
// 	"port_id": "3ccbeb73-63e5-4082-b90b-39859686de07",
// 	"project_id": "09522dddf25648d5bc30307cf3bf5f72",
// 	"qos_policy_id": null,
// 	"revision_number": 2,
// 	"router_id": "2efe00cc-f740-4b39-8364-b417340652b8",
// 	"status": "ACTIVE",
// 	"subnet_id": null,
// 	"tags": [],
// 	"updated_at": "2023-04-03T15:09:49Z"
//   }
