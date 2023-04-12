package models

import "time"

type Volume struct {
	Attachments      []VolumeAttachment `json:"attachments"`
	AvailabilityZone string            `json:"availability_zone"`
	Bootable         bool              `json:"bootable"`
	CreatedAt        time.Time         `json:"created_at"`
	Description      string            `json:"description"`
	Encrypted        bool              `json:"encrypted"`
	ID               string            `json:"id"`
	MultiAttach      bool              `json:"multiattach"`
	Name             string            `json:"name"`
	Host             string            `json:"host"`
	Properties       []Property        `json:"properties"`
	Size             int               `json:"size"`
	SnapshotID       string            `json:"snapshot_id"`
	SourceVolumeID   string            `json:"source_volid"`
	Status           string            `json:"status"`
	Type             string            `json:"type"`
	UpdatedAt        time.Time         `json:"updated_at"`
	AccountID        string            `json:"user_id"`
}

// {
// 	"attachments": [
// 	  {
// 		"id": "12a758e9-3cb0-4e0f-9482-a7e9712bc3c9",
// 		"attachment_id": "da6dc769-90d8-4c1c-8803-acec283057dd",
// 		"volume_id": "12a758e9-3cb0-4e0f-9482-a7e9712bc3c9",
// 		"server_id": "1a639b0f-30fa-47c7-9ba3-9512ce1583d6",
// 		"host_name": "compute.test.local",
// 		"device": "/dev/vdb",
// 		"attached_at": "2023-04-10T12:03:29.000000"
// 	  }
// 	],
// 	"availability_zone": "nova",
// 	"bootable": "false",
// 	"consistencygroup_id": null,
// 	"created_at": "2023-04-10T12:02:44.000000",
// 	"description": null,
// 	"encrypted": false,
// 	"id": "12a758e9-3cb0-4e0f-9482-a7e9712bc3c9",
// 	"migration_status": null,
// 	"multiattach": false,
// 	"name": "disk01",
// 	"os-vol-host-attr:host": "network.test.local@lvm#LVM",
// 	"os-vol-mig-status-attr:migstat": null,
// 	"os-vol-mig-status-attr:name_id": null,
// 	"os-vol-tenant-attr:tenant_id": "09522dddf25648d5bc30307cf3bf5f72",
// 	"properties": {},
// 	"replication_status": null,
// 	"size": 2,
// 	"snapshot_id": null,
// 	"source_volid": null,
// 	"status": "in-use",
// 	"type": "__DEFAULT__",
// 	"updated_at": "2023-04-10T12:03:30.000000",
// 	"user_id": "63cec20a0cd44e689d889fc164d179b7"
//   }
