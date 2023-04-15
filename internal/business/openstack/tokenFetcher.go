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
	req, err := generateReq(account)
	if err != nil {
		return nil, err
	}

	json_data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	cluster := f.fetcher.clusters[1]

	resp, err := f.fetcher.client.Post(cluster.URL, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}

	token := &models.Token{}

	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	token.Value = resp.Header.Get("X-Subject-Token")

	return token, nil
}

func generateReq(account *models.Account) (GetTokenRequest, error) {
	methods := [1]string{"password"}

	req := GetTokenRequest{
		Auth: &Identity{
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
	}

	return req, nil
}
