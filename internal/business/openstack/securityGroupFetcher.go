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

func (f *SecurityGroupFetcher) FetchByID(ctx context.Context, securityGroupId string) (*models.SecurityGroup, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchSecurityGroupURL := cluster.URL + ":9696" + "/v2.0/security-groups/" + securityGroupId

	req, err := http.NewRequest("GET", fetchSecurityGroupURL, nil)
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

	securityGroup := &models.SecurityGroup{}

	fetchSecurityGroupRes := &FetchSecurityGroupResponse{
		SecurityGroup: securityGroup,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchSecurityGroupRes); err != nil {
		return nil, err
	}

	return fetchSecurityGroupRes.SecurityGroup, nil
}

func (f *SecurityGroupFetcher) Create(ctx context.Context, securityGroup *models.SecurityGroup) error {
	reqData := f.generateCreateReq(securityGroup)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createSecurityGroupURL := cluster.URL + ":9696" + "/v2.0/security-groups"

	req, err := http.NewRequest("POST", createSecurityGroupURL, bytes.NewBuffer(json_data))
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

	createSecurityGroupRes := &CreateSecurityGroupResponse{
		SecurityGroup: securityGroup,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createSecurityGroupRes); err != nil {
		return err
	}

	return nil

}

func (f *SecurityGroupFetcher) Delete(ctx context.Context, securityGroupID string) error {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	deleteSecurityGroupURL := cluster.URL + ":9696" + "/v2.0/security-groups/" + securityGroupID

	req, err := http.NewRequest("DELETE", deleteSecurityGroupURL, nil)
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
		return errors.New("internal server error")
	}

	return nil
}

func (f *SecurityGroupFetcher) Update(ctx context.Context, securityGroupId string) {

}

func (f *SecurityGroupFetcher) generateCreateReq(securityGroup *models.SecurityGroup) *CreateSecurityGroupRequest {
	req := &CreateSecurityGroupRequest{
		SecurityGroup: &SecurityGroup{
			Name:        securityGroup.Name,
			Description: securityGroup.Description,
		},
	}

	return req
}
