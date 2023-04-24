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

func (f *FloatingIpFetcher) generateCreateReq(folatingIp *models.FloatingIp) *CreateFloatingIpRequest {
	req := &CreateFloatingIpRequest{
		FloatingIp: &FloatingIp{
			NetworkID: folatingIp.NetworkID,
		},
	}

	return req
}
