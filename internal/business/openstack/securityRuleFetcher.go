package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type SecurityRuleFetcher struct {
	fetcher *Fetcher
}

func (f *SecurityRuleFetcher) FetchByID(ctx context.Context, clusterId int, securityRuleId string) (*models.SecurityRule, error) {
	cluster := f.fetcher.clusters[clusterId-1]

	fetchSecurityRuleURL := cluster.URL + ":5000" + "/v3/securityRules/" + securityRuleId

	req, err := http.NewRequest("GET", fetchSecurityRuleURL, nil)
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

	securityRule := &models.SecurityRule{}

	fetchSecurityRuleRes := &FetchSecurityRuleResponse{
		SecurityRule: securityRule,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchSecurityRuleRes); err != nil {
		return nil, err
	}

	return fetchSecurityRuleRes.SecurityRule, nil
}

func (f *SecurityRuleFetcher) Create(ctx context.Context, clusterId int, securityRule *models.SecurityRule) error {
	reqData := f.generateCreateReq(securityRule)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	cluster := f.fetcher.clusters[clusterId-1]

	createSecurityRuleURL := cluster.URL + ":5000" + "/v3/securityRules"

	req, err := http.NewRequest("POST", createSecurityRuleURL, bytes.NewBuffer(json_data))
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

	createSecurityRuleRes := &CreateSecurityRuleResponse{
		SecurityRule: securityRule,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createSecurityRuleRes); err != nil {
		return err
	}

	return nil

}

func (f *SecurityRuleFetcher) Delete(ctx context.Context, clusterId int, securityRuleID string) error {
	cluster := f.fetcher.clusters[clusterId-1]

	deleteSecurityRuleURL := cluster.URL + ":5000" + "/v3/securityRules/" + securityRuleID

	req, err := http.NewRequest("DELETE", deleteSecurityRuleURL, nil)
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

func (f *SecurityRuleFetcher) Update(ctx context.Context, clusterId int) {

}

func (f *SecurityRuleFetcher) getAdminToken(ctx context.Context, clusterId int) (*models.Token, error) {

	token, err := f.fetcher.Token().Get(ctx, clusterId, f.fetcher.clusters[clusterId-1].Admin)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (f *SecurityRuleFetcher) generateCreateReq(securityRule *models.SecurityRule) *CreateSecurityRuleRequest {
	req := &CreateSecurityRuleRequest{
		SecurityRule: &SecurityRule{},
	}

	return req
}
