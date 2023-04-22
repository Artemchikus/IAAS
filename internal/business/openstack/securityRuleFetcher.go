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

func (f *SecurityRuleFetcher) FetchByID(ctx context.Context, securityRuleId string) (*models.SecurityRule, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchSecurityRuleURL := cluster.URL + ":9696" + "/v2.0/security-group-rules/" + securityRuleId

	req, err := http.NewRequest("GET", fetchSecurityRuleURL, nil)
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

	securityRule := &models.SecurityRule{}

	fetchSecurityRuleRes := &FetchSecurityRuleResponse{
		SecurityRule: securityRule,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchSecurityRuleRes); err != nil {
		return nil, err
	}

	return fetchSecurityRuleRes.SecurityRule, nil
}

func (f *SecurityRuleFetcher) Create(ctx context.Context, securityRule *models.SecurityRule) error {
	reqData := f.generateCreateReq(securityRule)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createSecurityRuleURL := cluster.URL + ":9696" + "/v2.0/security-group-rules"

	req, err := http.NewRequest("POST", createSecurityRuleURL, bytes.NewBuffer(json_data))
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

	createSecurityRuleRes := &CreateSecurityRuleResponse{
		SecurityRule: securityRule,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createSecurityRuleRes); err != nil {
		return err
	}

	return nil

}

func (f *SecurityRuleFetcher) Delete(ctx context.Context, securityRuleID string) error {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	deleteSecurityRuleURL := cluster.URL + ":9696" + "/v2.0/security-group-rules/" + securityRuleID

	req, err := http.NewRequest("DELETE", deleteSecurityRuleURL, nil)
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

func (f *SecurityRuleFetcher) Update(ctx context.Context, securityRuleId string) {

}

func (f *SecurityRuleFetcher) generateCreateReq(securityRule *models.SecurityRule) *CreateSecurityRuleRequest {
	req := &CreateSecurityRuleRequest{
		SecurityRule: &SecurityRule{
			Protocol:        securityRule.Protocol,
			PortRangeMax:    securityRule.PortRangeMax,
			PortRangeMin:    securityRule.PortRangeMin,
			RemoteIPPrefix:  securityRule.RemoteIpPrefix,
			Ethertype:       securityRule.Ethertype,
			Direction:       securityRule.Direction,
			SecurityGroupID: securityRule.SecurityGroupID,
		},
	}

	return req
}
