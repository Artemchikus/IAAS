package models

type Flavor struct {
	Disabled    bool     `json:"disabled"`
	Description string   `json:"description"`
	Disk        int      `json:"disk"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	IsPublic    bool     `json:"is_public"`
	Properties  []string `json:"properties"`
	RAM         int      `json:"ram"`
	RXTXFactor  float32  `json:"rxtx_factor"`
	Swap        string   `json:"swap"`
	VCPUs       int      `json:"vcpus"`
	Ephemeral   int      `json:"ephemeral"`
}

// {
// 	"OS-FLV-DISABLED:disabled": false,
// 	"OS-FLV-EXT-DATA:ephemeral": 0,
// 	"access_project_ids": null,
// 	"description": null,
// 	"disk": 2,
// 	"id": "0",
// 	"name": "m1.small",
// 	"os-flavor-access:is_public": true,
// 	"properties": {},
// 	"ram": 300,
// 	"rxtx_factor": 1.0,
// 	"swap": "",
// 	"vcpus": 1
//   }
