package openstack

import (
	"IAAS/internal/models"
	"context"
	"encoding/json"
	"net/http"
)

type PortFetcher struct {
	fetcher *Fetcher
}

func (f *PortFetcher) FetchByID(ctx context.Context, portId string) (*models.Port, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchPortURL := cluster.URL + ":9696" + "/v2.0/ports/" + portId

	req, err := http.NewRequest("GET", fetchPortURL, nil)
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

	port := &models.Port{}

	fetchPortRes := &FetchPortResponse{
		Port: port,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchPortRes); err != nil {
		return nil, err
	}

	return fetchPortRes.Port, nil
}

func (f *PortFetcher) FetchByRouterID(ctx context.Context, routerId string) ([]*models.Port, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchPortURL := cluster.URL + ":9696" + "/v2.0/ports?device_id=" + routerId

	req, err := http.NewRequest("GET", fetchPortURL, nil)
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

	ports := []*models.Port{}

	fetchPortsRes := &FetchPortsResponse{
		Ports: ports,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchPortsRes); err != nil {
		return nil, err
	}

	return fetchPortsRes.Ports, nil
}

func (f *PortFetcher) FetchByNetworkID(ctx context.Context, networkId string) ([]*models.Port, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchPortURL := cluster.URL + ":9696" + "/v2.0/ports?network_id=" + networkId

	req, err := http.NewRequest("GET", fetchPortURL, nil)
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

	ports := []*models.Port{}

	fetchPortsRes := &FetchPortsResponse{
		Ports: ports,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchPortsRes); err != nil {
		return nil, err
	}

	return fetchPortsRes.Ports, nil
}