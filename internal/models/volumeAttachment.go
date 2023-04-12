package models

import "time"

type VolumeAttachment struct {
	ID           string    `json:"id"`
	AttachmentID string    `json:"attachment_id"`
	VolumeID     string    `json:"volume_id"`
	ServerID     string    `json:"server_id"`
	HostName     string    `json:"host_name"`
	Device       string    `json:"device"`
	AttacedAt    time.Time `json:"attaced_at"`
}

// {
// 		"id": "12a758e9-3cb0-4e0f-9482-a7e9712bc3c9",
// 		"attachment_id": "da6dc769-90d8-4c1c-8803-acec283057dd",
// 		"volume_id": "12a758e9-3cb0-4e0f-9482-a7e9712bc3c9",
// 		"server_id": "1a639b0f-30fa-47c7-9ba3-9512ce1583d6",
// 		"host_name": "compute.test.local",
// 		"device": "/dev/vdb",
// 		"attached_at": "2023-04-10T12:03:29.000000"
// 	  }
