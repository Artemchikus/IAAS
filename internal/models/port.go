package models

import "time"

type Port struct {
	ID          string    `json:"id"`
	DeviceID    string    `json:"device_id"`
	DeviceOwner string    `json:"device_owner"`
	MacAddress  string    `json:"mac_address"`
	Name        string    `json:"name"`
	NetworkID   string    `json:"network_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// {
// 	"admin_state_up": true,
// 	"allowed_address_pairs": [],
// 	"binding_host_id": "compute.test.local",
// 	"binding_profile": {},
// 	"binding_vif_details": {
// 	  "port_filter": true,
// 	  "connectivity": "l2",
// 	  "bound_drivers": {
// 		"0": "ovn"
// 	  }
// 	},
// 	"binding_vif_type": "ovs",
// 	"binding_vnic_type": "normal",
// 	"created_at": "2023-04-03T14:50:38Z",
// 	"data_plane_status": null,
// 	"description": "",
// 	"device_id": "1a639b0f-30fa-47c7-9ba3-9512ce1583d6",
// 	"device_owner": "compute:nova",
// 	"device_profile": null,
// 	"dns_assignment": null,
// 	"dns_domain": null,
// 	"dns_name": null,
// 	"extra_dhcp_opts": [],
// 	"fixed_ips": [
// 	  {
// 		"subnet_id": "94a7c47e-517e-452b-8abc-9b9958663d10",
// 		"ip_address": "192.168.100.31"
// 	  }
// 	],
// 	"id": "3ccbeb73-63e5-4082-b90b-39859686de07",
// 	"ip_allocation": null,
// 	"mac_address": "fa:16:3e:83:0d:c1",
// 	"name": "",
// 	"network_id": "439e51bd-f03f-44c5-9006-72a68d401e4a",
// 	"numa_affinity_policy": null,
// 	"port_security_enabled": true,
// 	"project_id": "09522dddf25648d5bc30307cf3bf5f72",
// 	"propagate_uplink_status": null,
// 	"qos_network_policy_id": null,
// 	"qos_policy_id": null,
// 	"resource_request": null,
// 	"revision_number": 5,
// 	"security_group_ids": [
// 	  "1f2f9c20-92bd-4e27-b6b3-c8530cd0d49a"
// 	],
// 	"status": "DOWN",
// 	"tags": [],
// 	"trunk_details": null,
// 	"updated_at": "2023-04-10T10:28:01Z"
//   }
