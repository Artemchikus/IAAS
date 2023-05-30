package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type NetworkFetcher struct {
	fetcher *Fetcher
}

func (f *NetworkFetcher) FetchByID(ctx context.Context, networkId string) (*models.Network, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchNetworkURL := cluster.URL + ":9696" + "/v2.0/networks/" + networkId

	req, err := http.NewRequest("GET", fetchNetworkURL, nil)
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

	network := &models.Network{}

	fetchNetworkRes := &FetchNetworkResponse{
		Network: network,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchNetworkRes); err != nil {
		return nil, err
	}

	return fetchNetworkRes.Network, nil
}

func (f *NetworkFetcher) Create(ctx context.Context, network *models.Network) error {
	reqData := f.generateCreateReq(network)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createNetworkURL := cluster.URL + ":9696" + "/v2.0/networks"

	req, err := http.NewRequest("POST", createNetworkURL, bytes.NewBuffer(json_data))
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

	createNetworkRes := &CreateNetworkResponse{
		Network: network,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createNetworkRes); err != nil {
		return err
	}

	return nil

}

func (f *NetworkFetcher) Delete(ctx context.Context, networkID string) error {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	deleteNetworkURL := cluster.URL + ":9696" + "/v2.0/networks/" + networkID

	req, err := http.NewRequest("DELETE", deleteNetworkURL, nil)
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

func (f *NetworkFetcher) FetchAll(ctx context.Context) ([]*models.Network, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchNetworkURL := cluster.URL + ":9696" + "/v2.0/networks"

	req, err := http.NewRequest("GET", fetchNetworkURL, nil)
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

	networks := []*models.Network{}

	fetchNetworksRes := &FetchNetworksResponse{
		Networks: &networks,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchNetworksRes); err != nil {
		return nil, err
	}

	return networks, nil
}

func (f *NetworkFetcher) generateCreateReq(network *models.Network) *CreateNetworkRequest {
	req := &CreateNetworkRequest{
		Network: &Network{
			ProjectID:       network.ProjectID,
			Name:            network.Name,
			NetworkType:     network.NetworkType,
			External:        network.External,
			PhysicalNetwork: network.PhysicalNetwork,
			Description:     network.Description,
			MTU:             network.MTU,
		},
	}

	return req
}
