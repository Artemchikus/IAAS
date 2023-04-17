package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type FlavorFetcher struct {
	fetcher *Fetcher
}

func (f *FlavorFetcher) FetchByID(ctx context.Context, clusterId int, flavorId string) (*models.Flavor, error) {
	cluster := f.fetcher.clusters[clusterId-1]

	fetchFlavorURL := cluster.URL + ":8774" + "/v2.1/flavors/" + flavorId

	req, err := http.NewRequest("GET", fetchFlavorURL, nil)
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

	flavor := &models.Flavor{}

	fetchFlavorRes := &FetchFlavorResponse{
		Flavor: flavor,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchFlavorRes); err != nil {
		return nil, err
	}

	return fetchFlavorRes.Flavor, nil
}

func (f *FlavorFetcher) Create(ctx context.Context, clusterId int, flavor *models.Flavor) error {
	reqData := f.generateCreateReq(flavor)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	cluster := f.fetcher.clusters[clusterId-1]

	createFlavorURL := cluster.URL + ":8774" + "/v2.1/flavors"

	req, err := http.NewRequest("POST", createFlavorURL, bytes.NewBuffer(json_data))
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

	createFlavorRes := &CreateFlavorResponse{
		Flavor: flavor,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createFlavorRes); err != nil {
		return err
	}

	return nil
}
func (f *FlavorFetcher) Delete(ctx context.Context, clusterId int, flavorId string) error {
	cluster := f.fetcher.clusters[clusterId-1]

	deleteFlavorURL := cluster.URL + ":8774" + "/v2.1/flavors/" + flavorId

	req, err := http.NewRequest("DELETE", deleteFlavorURL, nil)
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

	if resp.StatusCode != 202 {
		return errors.New("internal server error")
	}

	return nil
}

func (f *FlavorFetcher) Update(ctx context.Context, clusterId int) {

}

func (f *FlavorFetcher) getAdminToken(ctx context.Context, clusterId int) (*models.Token, error) {

	token, err := f.fetcher.Token().Get(ctx, clusterId, f.fetcher.clusters[clusterId-1].Admin)
	if err != nil {
		return nil, err
	}

	return token, nil
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
