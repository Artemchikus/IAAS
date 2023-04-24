package openstack

import (
	"IAAS/internal/models"
	"context"
)

type ServerFetcher struct {
	fetcher *Fetcher
}

func (f *ServerFetcher) FetchByID(ctx context.Context, serverId string) (*models.Server, error) {
	return nil, nil
}

func (f *ServerFetcher) Create(ctx context.Context, server *models.Server) error {
	return nil
}

func (f *ServerFetcher) Delete(ctx context.Context, serverId string) error {
	return nil
}

func (f *ServerFetcher) AddVolume(ctx context.Context, serverId, volumeId string) error {
	// 	REQ: curl -g -i -X POST http://controller.test.local:8774/v2.1/servers/fb8af66c-d618-4217-9742-30017b78a636/os-volume_attachments -H "Content-Type: application/json" -H "OpenStack-API-Version: compute 2.89" -H "User-Agent: openstacksdk/0.101.0 keystoneauth1/5.0.0 python-requests/2.25.1 CPython/3.9.16" -H "X-Auth-Token: {SHA256}8307b83f055b813bc96afb1e02bf61d06cb22c0c7b9261528fd826fb73c5ab31" -H "X-OpenStack-Nova-API-Version: 2.89" -d '{"volumeAttachment": {"volumeId": "a14efe17-a246-4495-b330-73b1c93162a0", "device": null}}'
	// http://controller.test.local:8774 "POST /v2.1/servers/fb8af66c-d618-4217-9742-30017b78a636/os-volume_attachments HTTP/1.1" 200 239
	return nil
}

func (f *ServerFetcher) AddFloatingIp(ctx context.Context, serverId, floatingIpId string) error {
	// 	REQ: curl -g -i -X PUT http://controller.test.local:9696/v2.0/floatingips/039a4f3d-5a52-4127-8c25-41b3795ec980 -H "Content-Type: application/json" -H "User-Agent: openstacksdk/0.101.0 keystoneauth1/5.0.0 python-requests/2.25.1 CPython/3.9.16" -H "X-Auth-Token: {SHA256}4d64e3bf32e1a3eeb270a6355229a52272f7b142538d4963bd140ec3efa557b0" -d '{"floatingip": {"port_id": "ba8c8d6c-2ee9-4035-adac-c3c5b7bad3c5"}}'
	// http://controller.test.local:9696 "PUT /v2.0/floatingips/039a4f3d-5a52-4127-8c25-41b3795ec980 HTTP/1.1" 200 838
	return nil
}
