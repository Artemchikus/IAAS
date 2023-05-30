package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type FloatingIpFetcher struct {
	fetcher *Fetcher
}

func (f *FloatingIpFetcher) FetchByID(ctx context.Context, folatingIpId string) (*models.FloatingIp, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchFloatingIpURL := cluster.URL + ":9696" + "/v2.0/floatingips/" + folatingIpId

	req, err := http.NewRequest("GET", fetchFloatingIpURL, nil)
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

	folatingIp := &models.FloatingIp{}

	fetchFloatingIpRes := &FetchFloatingIpResponse{
		FloatingIp: folatingIp,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchFloatingIpRes); err != nil {
		return nil, err
	}

	return fetchFloatingIpRes.FloatingIp, nil
}

func (f *FloatingIpFetcher) Create(ctx context.Context, folatingIp *models.FloatingIp) error {
	reqData := f.generateCreateReq(folatingIp)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createFloatingIpURL := cluster.URL + ":9696" + "/v2.0/floatingips"

	req, err := http.NewRequest("POST", createFloatingIpURL, bytes.NewBuffer(json_data))
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

	if resp.StatusCode != 201 {
		return handleErrorResponse(resp)
	}

	createFloatingIpRes := &CreateFloatingIpResponse{
		FloatingIp: folatingIp,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createFloatingIpRes); err != nil {
		return err
	}

	return nil

}

func (f *FloatingIpFetcher) Delete(ctx context.Context, folatingIpID string) error {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	deleteFloatingIpURL := cluster.URL + ":9696" + "/v2.0/floatingips/" + folatingIpID

	req, err := http.NewRequest("DELETE", deleteFloatingIpURL, nil)
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

func (f *FloatingIpFetcher) AddToPort(ctx context.Context, floatingIpID, portID string) error {
	// 	REQ: curl -g -i -X PUT http://controller.test.local:9696/v2.0/floatingips/039a4f3d-5a52-4127-8c25-41b3795ec980 -H "Content-Type: application/json" -H "User-Agent: openstacksdk/0.101.0 keystoneauth1/5.0.0 python-requests/2.25.1 CPython/3.9.16" -H "X-Auth-Token: {SHA256}4d64e3bf32e1a3eeb270a6355229a52272f7b142538d4963bd140ec3efa557b0" -d '{"floatingip": {"port_id": "ba8c8d6c-2ee9-4035-adac-c3c5b7bad3c5"}}'
	// http://controller.test.local:9696 "PUT /v2.0/floatingips/039a4f3d-5a52-4127-8c25-41b3795ec980 HTTP/1.1" 200 838
	ip := &AddIpToPort{
		PortID: portID,
	}

	reqData := &AddIpToPortRequest{
		FloatingIp: ip,
	}

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createServerURL := cluster.URL + ":9696" + "/v2.0/floatingips/" + floatingIpID

	req, err := http.NewRequest("PUT", createServerURL, bytes.NewBuffer(json_data))
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

func (f *FloatingIpFetcher) FetchAll(ctx context.Context) ([]*models.FloatingIp, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchFloatingIpURL := cluster.URL + ":9696" + "/v2.0/floatingips"

	req, err := http.NewRequest("GET", fetchFloatingIpURL, nil)
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

	floatingips := []*models.FloatingIp{}

	fetchFloatingIpsRes := &FetchFloatingIpsResponse{
		FloatingIps: &floatingips,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchFloatingIpsRes); err != nil {
		return nil, err
	}

	return floatingips, nil
}

func (f *FloatingIpFetcher) generateCreateReq(floatingIp *models.FloatingIp) *CreateFloatingIpRequest {
	req := &CreateFloatingIpRequest{
		FloatingIp: &FloatingIp{
			NetworkID:   floatingIp.NetworkID,
			Description: floatingIp.Description,
		},
	}

	return req
}
