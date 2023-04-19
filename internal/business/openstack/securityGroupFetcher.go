package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type SecurityGroupFetcher struct {
	fetcher *Fetcher
}

func (f *SecurityGroupFetcher) FetchByID(ctx context.Context, clusterId int, securityGroupId string) (*models.SecurityGroup, error) {
	cluster := f.fetcher.clusters[clusterId-1]

	fetchSecurityGroupURL := cluster.URL + ":5000" + "/v3/securityGroups/" + securityGroupId

	req, err := http.NewRequest("GET", fetchSecurityGroupURL, nil)
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

	securityGroup := &models.SecurityGroup{}

	fetchSecurityGroupRes := &FetchSecurityGroupResponse{
		SecurityGroup: securityGroup,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchSecurityGroupRes); err != nil {
		return nil, err
	}

	return fetchSecurityGroupRes.SecurityGroup, nil
}

func (f *SecurityGroupFetcher) Create(ctx context.Context, clusterId int, securityGroup *models.SecurityGroup) error {
	reqData := f.generateCreateReq(securityGroup)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	cluster := f.fetcher.clusters[clusterId-1]

	createSecurityGroupURL := cluster.URL + ":5000" + "/v3/securityGroups"

	req, err := http.NewRequest("POST", createSecurityGroupURL, bytes.NewBuffer(json_data))
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

	createSecurityGroupRes := &CreateSecurityGroupResponse{
		SecurityGroup: securityGroup,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createSecurityGroupRes); err != nil {
		return err
	}

	return nil

}

func (f *SecurityGroupFetcher) Delete(ctx context.Context, clusterId int, securityGroupID string) error {
	cluster := f.fetcher.clusters[clusterId-1]

	deleteSecurityGroupURL := cluster.URL + ":5000" + "/v3/securityGroups/" + securityGroupID

	req, err := http.NewRequest("DELETE", deleteSecurityGroupURL, nil)
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

func (f *SecurityGroupFetcher) Update(ctx context.Context, clusterId int) {

}

func (f *SecurityGroupFetcher) getAdminToken(ctx context.Context, clusterId int) (*models.Token, error) {

	token, err := f.fetcher.Token().Get(ctx, clusterId, f.fetcher.clusters[clusterId-1].Admin)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (f *SecurityGroupFetcher) generateCreateReq(securityGroup *models.SecurityGroup) *CreateSecurityGroupRequest {
	req := &CreateSecurityGroupRequest{
		SecurityGroup: &SecurityGroup{},
	}

	return req
}
