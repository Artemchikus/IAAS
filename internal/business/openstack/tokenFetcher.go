package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
)

type TokenFetcher struct {
	fetcher *Fetcher
}

func (f *TokenFetcher) Get(ctx context.Context, clusterId int, account *models.Account) (*models.Token, error) {
	req := f.generateGetReq(account)

	json_data, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	cluster := f.fetcher.clusters[clusterId-1]

	getTokenURL := cluster.URL + ":5000" + "/v3/auth/tokens"

	resp, err := f.fetcher.client.Post(getTokenURL, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	token := &models.Token{}
	tokenRes := &GetTokenResponse{
		Token: token,
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return nil, err
	}

	token.Value = resp.Header.Get("X-Subject-Token")

	return token, nil
}

func (f *TokenFetcher) generateGetReq(account *models.Account) *GetTokenRequest {
	methods := [1]string{"password"}

	req := &GetTokenRequest{
		Auth: &Auth{
			Identity: &Identity{
				Methods: methods[:],
				Password: &Password{
					User: &User{
						Name: account.Name,
						Domain: &Domain{
							ID: "default",
						},
						Password: account.Password,
					},
				},
			},
			Scope: &Scope{
				Project: &GetTokenProject{
					ID: account.ProjectID,
					Domain: &Domain{
						ID: "default",
					},
				},
			},
		},
	}

	return req
}
