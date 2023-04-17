package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type SubnetFetcher struct {
	fetcher *Fetcher
}

func (f *SubnetFetcher) FetchByID(ctx context.Context, clusterId int, subnetId string) (*models.Subnet, error) {
	cluster := f.fetcher.clusters[clusterId-1]

	fetchSubnetURL := cluster.URL + ":9696" + "/v2.0/subnets/" + subnetId

	req, err := http.NewRequest("GET", fetchSubnetURL, nil)
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

	subnet := &models.Subnet{}

	fetchSubnetRes := &FetchSubnetResponse{
		Subnet: subnet,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchSubnetRes); err != nil {
		return nil, err
	}

	return subnet, nil
}

func (f *SubnetFetcher) Create(ctx context.Context, clusterId int, subnet *models.Subnet) error {
	reqData := f.generateCreateReq(subnet)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	cluster := f.fetcher.clusters[clusterId-1]

	createSubnetURL := cluster.URL + ":9696" + "/v2.0/subnets"

	req, err := http.NewRequest("POST", createSubnetURL, bytes.NewBuffer(json_data))
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

	createSubnetRes := &FetchSubnetResponse{
		Subnet: subnet,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createSubnetRes); err != nil {
		return err
	}

	return nil

}

func (f *SubnetFetcher) Delete(ctx context.Context, clusterId int, subnetID string) error {
	cluster := f.fetcher.clusters[clusterId-1]

	deleteSubnetURL := cluster.URL + ":9696" + "/v2.0/subnets/" + subnetID

	req, err := http.NewRequest("DELETE", deleteSubnetURL, nil)
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

func (f *SubnetFetcher) Update(ctx context.Context, clusterId int) {

}

func (f *SubnetFetcher) getAdminToken(ctx context.Context, clusterId int) (*models.Token, error) {

	token, err := f.fetcher.Token().Get(ctx, clusterId, f.fetcher.clusters[clusterId-1].Admin)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (f *SubnetFetcher) generateCreateReq(subnet *models.Subnet) *CreateSubnetRequest {
	allocationPools := []*AllocationPool{}

	for _, pool := range subnet.AllocationPools {
		newPool := &AllocationPool{
			Start: pool.Start,
			End:   pool.End,
		}
		allocationPools = append(allocationPools, newPool)
	}

	req := &CreateSubnetRequest{
		Subnet: &Subnet{
			CIDR:            subnet.CIDR,
			Name:            subnet.Name,
			EnableDHCP:      subnet.EnableDHCP,
			NetworkID:       subnet.NetworkID,
			AllocationPools: allocationPools,
			IpVersion:       subnet.IpVersion,
			GatewayIp:       subnet.GatewayIp,
		},
	}

	return req
}