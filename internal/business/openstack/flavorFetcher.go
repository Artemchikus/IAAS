package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type FlavorFetcher struct {
	fetcher *Fetcher
}

func (f *FlavorFetcher) FetchByID(ctx context.Context, flavorId string) (*models.Flavor, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchFlavorURL := cluster.URL + ":8774" + "/v2.1/flavors/" + flavorId

	req, err := http.NewRequest("GET", fetchFlavorURL, nil)
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

	flavor := &models.Flavor{}

	fetchFlavorRes := &FetchFlavorResponse{
		Flavor: flavor,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchFlavorRes); err != nil {
		return nil, err
	}

	return fetchFlavorRes.Flavor, nil
}

func (f *FlavorFetcher) Create(ctx context.Context, flavor *models.Flavor) error {
	reqData := f.generateCreateReq(flavor)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createFlavorURL := cluster.URL + ":8774" + "/v2.1/flavors"

	req, err := http.NewRequest("POST", createFlavorURL, bytes.NewBuffer(json_data))
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

	createFlavorRes := &CreateFlavorResponse{
		Flavor: flavor,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createFlavorRes); err != nil {
		return err
	}

	return nil
}

func (f *FlavorFetcher) Delete(ctx context.Context, flavorId string) error {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	deleteFlavorURL := cluster.URL + ":8774" + "/v2.1/flavors/" + flavorId

	req, err := http.NewRequest("DELETE", deleteFlavorURL, nil)
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

	if resp.StatusCode != 202 {
		return handleErrorResponse(resp)
	}

	return nil
}

func (f *FlavorFetcher) FetchAll(ctx context.Context) ([]*models.Flavor, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchFlavorURL := cluster.URL + ":8774" + "/v2.1/flavors/detail"

	req, err := http.NewRequest("GET", fetchFlavorURL, nil)
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

	flavors := []*models.Flavor{}

	fetchFlavorsRes := &FetchFlavorsResponse{
		Flavors: &flavors,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchFlavorsRes); err != nil {
		return nil, err
	}

	return flavors, nil
}

func (f *FlavorFetcher) generateCreateReq(flavor *models.Flavor) *CreateFlavorRequest {
	req := &CreateFlavorRequest{
		Flavor: &Flavor{
			Name:       flavor.Name,
			VCPUs:      flavor.VCPUs,
			Disk:       flavor.Disk,
			RAM:        flavor.RAM,
			IsPublic:   flavor.IsPublic,
			RXTXFactor: flavor.RXTXFactor,
			Swap:       flavor.Swap,
			Ephemeral:  flavor.Ephemeral,
		},
	}

	return req
}
