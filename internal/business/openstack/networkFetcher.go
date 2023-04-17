package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type NetworkFetcher struct {
	fetcher *Fetcher
}

func (f *NetworkFetcher) FetchByID(ctx context.Context, clusterId int, networkId string) (*models.Network, error) {
	cluster := f.fetcher.clusters[clusterId-1]

	fetchNetworkURL := cluster.URL + ":9696" + "/v2.0/networks/" + networkId

	req, err := http.NewRequest("GET", fetchNetworkURL, nil)
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

	network := &models.Network{}

	fetchNetworkRes := &FetchNetworkResponse{
		Network: network,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchNetworkRes); err != nil {
		return nil, err
	}

	return fetchNetworkRes.Network, nil
}

func (f *NetworkFetcher) Create(ctx context.Context, clusterId int, network *models.Network) error {
	reqData := f.generateCreateReq(network)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	cluster := f.fetcher.clusters[clusterId-1]

	createNetworkURL := cluster.URL + ":9696" + "/v2.0/networks"

	req, err := http.NewRequest("POST", createNetworkURL, bytes.NewBuffer(json_data))
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

	createNetworkRes := &CreateNetworkResponse{
		Network: network,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createNetworkRes); err != nil {
		return err
	}

	return nil

}

func (f *NetworkFetcher) Delete(ctx context.Context, clusterId int, networkID string) error {
	cluster := f.fetcher.clusters[clusterId-1]

	deleteNetworkURL := cluster.URL + ":9696" + "/v2.0/networks/" + networkID

	req, err := http.NewRequest("DELETE", deleteNetworkURL, nil)
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

func (f *NetworkFetcher) Update(ctx context.Context, clusterId int) {

}

func (f *NetworkFetcher) getAdminToken(ctx context.Context, clusterId int) (*models.Token, error) {

	token, err := f.fetcher.Token().Get(ctx, clusterId, f.fetcher.clusters[clusterId-1].Admin)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (f *NetworkFetcher) generateCreateReq(network *models.Network) *CreateNetworkRequest {
	req := &CreateNetworkRequest{
		Network: &Network{
			Name:            network.Name,
			NetworkType:     network.NetworkType,
			AdminStateUp:    network.AdminStateUp,
			External:        network.External,
			PhysicalNetwork: network.PhysicalNetwork,
		},
	}

	return req
}
