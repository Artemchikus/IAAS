package models

import "time"

type Image struct {
	ID              string    `json:"id"`
	FileData        []byte    `json:"file_data"`
	Name            string    `json:"name"`
	DiskFormat      string    `json:"disk_format"`
	ContainerFormat string    `json:"container_format"`
	OwnerID         string    `json:"owner"`
	Visibility      string    `json:"visibility"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewImage(Name, DiskFormat, ContainerFormat, OwnerID, Visibility string) *Image {
	return &Image{
		Name:            Name,
		DiskFormat:      DiskFormat,
		ContainerFormat: ContainerFormat,
		Visibility:      Visibility,
	}
}

// {
//    "owner_specified.openstack.md5":"",
//    "owner_specified.openstack.sha256":"",
//    "owner_specified.openstack.object":"images/Cirros",
//    "name":"Cirros",
//    "disk_format":"qcow2",
//    "container_format":"bare",
//    "visibility":"public",
//    "size":null,
//    "virtual_size":null,
//    "status":"queued",
//    "checksum":null,
//    "protected":false,
//    "min_ram":0,
//    "min_disk":0,
//    "owner":"a0f8bb68921b4ddea63c52ef1b4e7e74",
//    "os_hidden":false,
//    "os_hash_algo":null,
//    "os_hash_value":null,
//    "id":"5e9b4878-a158-4c0d-8157-78bbb08b385a",
//    "created_at":"2023-04-17T11:08:29Z",
//    "updated_at":"2023-04-17T11:08:29Z",
//    "tags":[

//    ],
//    "self":"/v2/images/5e9b4878-a158-4c0d-8157-78bbb08b385a",
//    "file":"/v2/images/5e9b4878-a158-4c0d-8157-78bbb08b385a/file",
//    "schema":"/v2/schemas/image"
// }
