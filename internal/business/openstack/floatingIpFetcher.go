package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type FloatingIpFetcher struct {
	fetcher *Fetcher
}

func (f *FloatingIpFetcher) FetchByID(ctx context.Context, clusterId int, folatingIpId string) (*models.FloatingIp, error) {
	cluster := f.fetcher.clusters[clusterId-1]

	fetchFloatingIpURL := cluster.URL + ":9696" + "/v2.0/floatingips/" + folatingIpId

	req, err := http.NewRequest("GET", fetchFloatingIpURL, nil)
	if err != nil {
		return nil, err
	}

	token, err := f.getAdminToken(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Token", token.Value)

	resp, err := f.fetcher.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	folatingIp := &models.FloatingIp{}

	fetchFloatingIpRes := &FetchFloatingIpResponse{
		FloatingIp: folatingIp,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchFloatingIpRes); err != nil {
		return nil, err
	}

	return fetchFloatingIpRes.FloatingIp, nil
}

func (f *FloatingIpFetcher) Create(ctx context.Context, clusterId int, folatingIp *models.FloatingIp) error {
	reqData := f.generateCreateReq(folatingIp)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	cluster := f.fetcher.clusters[clusterId-1]

	createFloatingIpURL := cluster.URL + ":9696" + "/v2.0/floatingips"

	req, err := http.NewRequest("POST", createFloatingIpURL, bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}

	token, err := f.getAdminToken(ctx, clusterId)
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

	createFloatingIpRes := &CreateFloatingIpResponse{
		FloatingIp: folatingIp,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createFloatingIpRes); err != nil {
		return err
	}

	return nil

}

func (f *FloatingIpFetcher) Delete(ctx context.Context, clusterId int, folatingIpID string) error {
	cluster := f.fetcher.clusters[clusterId-1]

	deleteFloatingIpURL := cluster.URL + ":9696" + "/v2.0/floatingips/" + folatingIpID

	req, err := http.NewRequest("DELETE", deleteFloatingIpURL, nil)
	if err != nil {
		return err
	}

	token, err := f.getAdminToken(ctx, clusterId)
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
		return errors.New("internal server error")
	}

	return nil
}

func (f *FloatingIpFetcher) Update(ctx context.Context, clusterId int) {

}

func (f *FloatingIpFetcher) getAdminToken(ctx context.Context, clusterId int) (*models.Token, error) {

	token, err := f.fetcher.Token().Get(ctx, clusterId, f.fetcher.clusters[clusterId-1].Admin)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (f *FloatingIpFetcher) generateCreateReq(folatingIp *models.FloatingIp) *CreateFloatingIpRequest {
	req := &CreateFloatingIpRequest{
		FloatingIp: &FloatingIp{
			NetworkID: folatingIp.NetworkID,
		},
	}

	return req
}
