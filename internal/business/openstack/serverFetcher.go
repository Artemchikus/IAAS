package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ServerFetcher struct {
	fetcher *Fetcher
}

func (f *ServerFetcher) FetchByID(ctx context.Context, serverId string) (*models.Server, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchServerURL := cluster.URL + ":8774" + "/v2.1/servers/" + serverId

	req, err := http.NewRequest("GET", fetchServerURL, nil)
	if err != nil {
		return nil, err
	}

	token := getTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", token.Value)

	resp, err := f.fetcher.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, handleErrorResponse(resp)
	}

	serverResp := &FetchedServer{}

	fetchServerRes := &FetchServerResponse{
		Server: serverResp,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchServerRes); err != nil {
		return nil, err
	}

	var publicIp, privateIp, privateNetwork, volumeId string

	for k, v := range serverResp.Addresses {
		net := v

		privateNetwork = k

		for _, addr := range net {
			if addr.Type == "fixed" {
				privateIp = addr.Address
				continue
			}
			publicIp = addr.Address
		}
	}

	launchedAt, err := time.Parse("2006-01-02T15:04:05.000000", serverResp.LaunchedAt)
	if err != nil {
		return nil, err
	}

	if len(serverResp.Volumes) > 0 {
		volumeId = serverResp.Volumes[0].ID
	}

	server := &models.Server{
		HypervisorHostname: serverResp.HypervisorHostname,
		InstanceName:       serverResp.InstanceName,
		VMState:            serverResp.VMState,
		CreatedAt:          serverResp.CreatedAt,
		LaunchedAt:         launchedAt,
		UpdatedAt:          serverResp.UpdatedAt,
		PrivateIp:          privateIp,
		PublicIp:           publicIp,
		ID:                 serverResp.ID,
		ImageID:            serverResp.Image.ID,
		KeyID:              serverResp.Key,
		FlavorID:           serverResp.Flavor.ID,
		Name:               serverResp.Name,
		ProjectID:          serverResp.TenantID,
		SecurityGroupID:    serverResp.SecurityGroups[0].Name,
		Status:             serverResp.Status,
		UserID:             serverResp.UserID,
		VolumeAttachedID:   volumeId,
		PrivateNetworkID:   privateNetwork,
	}

	return server, nil
}

func (f *ServerFetcher) Create(ctx context.Context, server *models.Server) error {
	reqData := f.generateCreateReq(server)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createServerURL := cluster.URL + ":8774" + "/v2.1/servers"

	req, err := http.NewRequest("POST", createServerURL, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	token := getTokenFromContext(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", token.Value)
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.fetcher.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return handleErrorResponse(resp)
	}

	serverRes := map[string]interface{}{}

	createServerRes := &CreateServerResponse{
		Server: serverRes,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createServerRes); err != nil {
		return err
	}

	server.ID = serverRes["id"].(string)

	return nil
}

func (f *ServerFetcher) Delete(ctx context.Context, serverId string) error {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	deleteServerURL := cluster.URL + ":8774" + "/v2.1/servers/" + serverId

	req, err := http.NewRequest("DELETE", deleteServerURL, nil)
	if err != nil {
		return err
	}

	token := getTokenFromContext(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", token.Value)

	resp, err := f.fetcher.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return handleErrorResponse(resp)
	}

	return nil
}

func (f *ServerFetcher) Start(ctx context.Context, serverId string) error {
	reqData := &StartServerRequest{}

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createServerURL := cluster.URL + ":8774" + "/v2.1/servers/" + serverId + "/action"

	req, err := http.NewRequest("POST", createServerURL, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	token := getTokenFromContext(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", token.Value)
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.fetcher.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("resp: %v\n", resp)

	if resp.StatusCode != 202 {
		return handleErrorResponse(resp)
	}

	return nil
}

func (f *ServerFetcher) Stop(ctx context.Context, serverId string) error {
	reqData := &StopServerRequest{}

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createServerURL := cluster.URL + ":8774" + "/v2.1/servers/" + serverId + "/action"

	req, err := http.NewRequest("POST", createServerURL, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	token := getTokenFromContext(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", token.Value)
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.fetcher.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		return handleErrorResponse(resp)
	}

	return nil
}

func (f *ServerFetcher) AttachVolume(ctx context.Context, serverId, volumeId string) error {
	// 	REQ: curl -g -i -X POST http://controller.test.local:8774/v2.1/servers/fb8af66c-d618-4217-9742-30017b78a636/os-volume_attachments -H "Content-Type: application/json" -H "OpenStack-API-Version: compute 2.89" -H "User-Agent: openstacksdk/0.101.0 keystoneauth1/5.0.0 python-requests/2.25.1 CPython/3.9.16" -H "X-Auth-Token: {SHA256}8307b83f055b813bc96afb1e02bf61d06cb22c0c7b9261528fd826fb73c5ab31" -H "X-OpenStack-Nova-API-Version: 2.89" -d '{"volumeAttachment": {"volumeId": "a14efe17-a246-4495-b330-73b1c93162a0", "device": null}}'
	// http://controller.test.local:8774 "POST /v2.1/servers/fb8af66c-d618-4217-9742-30017b78a636/os-volume_attachments HTTP/1.1" 200 239
	volumeAttach := &VolumeAttachment{
		VolumeId: volumeId,
	}

	reqData := &AddVolumeToServerRequest{
		VolumeAttachment: volumeAttach,
	}

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createServerURL := cluster.URL + ":8774" + "/v2.1/servers/" + serverId + "/os-volume_attachments"

	req, err := http.NewRequest("POST", createServerURL, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	token := getTokenFromContext(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", token.Value)
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.fetcher.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return handleErrorResponse(resp)
	}

	return nil
}

func (f *ServerFetcher) FetchAll(ctx context.Context) ([]*models.Server, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchServerURL := cluster.URL + ":8774" + "/v2.1/servers/detail"

	req, err := http.NewRequest("GET", fetchServerURL, nil)
	if err != nil {
		return nil, err
	}

	token := getTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", token.Value)

	resp, err := f.fetcher.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, handleErrorResponse(resp)
	}

	fetchedServers := []*FetchedServer{}

	fetchServersRes := &FetchServersResponse{
		Servers: &fetchedServers,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchServersRes); err != nil {
		return nil, err
	}

	fmt.Printf("fetchedServers: %v\n", fetchedServers[0])

	servers := []*models.Server{}

	var publicIp, privateIp, privateNetwork, volumeId string

	for _, server := range fetchedServers {
		for k, v := range server.Addresses {
			net := v

			privateNetwork = k

			for _, addr := range net {
				if addr.Type == "fixed" {
					privateIp = addr.Address
					continue
				}
				publicIp = addr.Address
			}
		}

		launchedAt, err := time.Parse("2006-01-02T15:04:05.000000", server.LaunchedAt)
		if err != nil {
			return nil, err
		}

		if len(server.Volumes) > 0 {
			volumeId = server.Volumes[0].ID
		}

		s := &models.Server{
			HypervisorHostname: server.HypervisorHostname,
			InstanceName:       server.InstanceName,
			VMState:            server.VMState,
			CreatedAt:          server.CreatedAt,
			LaunchedAt:         launchedAt,
			UpdatedAt:          server.UpdatedAt,
			PrivateIp:          privateIp,
			PublicIp:           publicIp,
			ID:                 server.ID,
			ImageID:            server.Image.ID,
			KeyID:              server.Key,
			FlavorID:           server.Flavor.ID,
			Name:               server.Name,
			ProjectID:          server.TenantID,
			SecurityGroupID:    server.SecurityGroups[0].Name,
			Status:             server.Status,
			UserID:             server.UserID,
			VolumeAttachedID:   volumeId,
			PrivateNetworkID:   privateNetwork,
		}

		servers = append(servers, s)
	}

	return servers, nil
}

func (f *ServerFetcher) generateCreateReq(server *models.Server) *CreateServerRequest {
	usg := &ServerSecurityGroupID{
		ID: server.SecurityGroupID,
	}

	unet := &ServerNetworkID{
		ID: server.PrivateNetworkID,
	}

	sgs := []*ServerSecurityGroupID{}
	sgs = append(sgs, usg)

	nets := []*ServerNetworkID{}
	nets = append(nets, unet)

	req := &CreateServerRequest{
		Server: &Server{
			FlavorID:       server.FlavorID,
			ImageID:        server.ImageID,
			KeyID:          server.KeyID,
			Name:           server.Name,
			SecurityGroups: sgs,
			Networks:       nets,
		},
	}

	return req
}
