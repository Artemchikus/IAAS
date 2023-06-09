package models

import "time"

type Server struct {
	HypervisorHostname string    `json:"hypervisor_hostname"`
	InstanceName       string    `json:"instance_name"`
	VMState            string    `json:"vm_state"`
	CreatedAt          time.Time `json:"created_at"`
	LaunchedAt         time.Time `json:"launched_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	PrivateIp          string    `json:"private_ip"`
	PublicIp           string    `json:"public_ip"`
	ID                 string    `json:"id"`
	ImageID            string    `json:"image_id"`
	KeyID              string    `json:"key_name"`
	FlavorID           string    `json:"flavor"`
	Name               string    `json:"name"`
	ProjectID          string    `json:"project_id"`
	SecurityGroupID    string    `json:"security_groups"`
	Status             string    `json:"status"`
	UserID             string    `json:"user_id"`
	VolumeAttachedID   string    `json:"volume_attached"`
	PrivateNetworkID   string    `json:"private_network_id"`
}

func NewServer(ImageID, KeyID, Name, SecurityGroupID, PrivateNetworkID, FlavorID string) *Server {
	return &Server{
		ImageID:          ImageID,
		KeyID:            KeyID,
		Name:             Name,
		SecurityGroupID:  SecurityGroupID,
		PrivateNetworkID: PrivateNetworkID,
		FlavorID:         FlavorID,
	}
}

// {
// 	"OS-DCF:diskConfig": "MANUAL",
// 	"OS-EXT-AZ:availability_zone": "nova",
// 	"OS-EXT-SRV-ATTR:host": "compute.test.local",
// 	"OS-EXT-SRV-ATTR:hypervisor_hostname": "compute.test.local",
// 	"OS-EXT-SRV-ATTR:instance_name": "instance-00000004",
// 	"OS-EXT-STS:power_state": 4,
// 	"OS-EXT-STS:task_state": null,
// 	"OS-EXT-STS:vm_state": "stopped",
// 	"OS-SRV-USG:launched_at": "2023-04-03T14:54:34.000000",
// 	"OS-SRV-USG:terminated_at": null,
// 	"accessIPv4": "",
// 	"accessIPv6": "",
// 	"addresses": {
// 	  "private": [
// 		"192.168.100.31",
// 		"192.168.122.213"
// 	  ]
// 	},
// 	"config_drive": "",
// 	"created": "2023-04-03T14:50:37Z",
// 	"flavor": "m1.small (0)",
// 	"hostId": "971721fcfe9b3aa38e8f0ffc15f9bc14e5d18e87a68fae8f98d37128",
// 	"id": "1a639b0f-30fa-47c7-9ba3-9512ce1583d6",
// 	"image": "Cirros (8409386d-bf2b-4433-8bbc-417a5f8594e6)",
// 	"key_name": "mykey",
// 	"name": "Cirros",
// 	"project_id": "09522dddf25648d5bc30307cf3bf5f72",
// 	"properties": {},
// 	"security_groups": [
// 	  {
// 		"name": "secgroup01"
// 	  }
// 	],
// 	"status": "SHUTOFF",
// 	"updated": "2023-04-10T09:48:51Z",
// 	"user_id": "63cec20a0cd44e689d889fc164d179b7",
// 	"volumes_attached": []
//   }
