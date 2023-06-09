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

func (f *TokenFetcher) Get(ctx context.Context, clusteUser *models.ClusterUser) (*models.Token, error) {
	clusterId := getClusterIDFromContext(ctx)

	req := f.generateGetReq(clusteUser)

	json_data, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	cluster := f.fetcher.clusters[clusterId]

	getTokenURL := cluster.URL + ":5000" + "/v3/auth/tokens"

	resp, err := f.fetcher.client.Post(getTokenURL, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, handleErrorResponse(resp)
	}

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

func (f *TokenFetcher) GetAdmin(ctx context.Context) (*models.Token, error) {
	clusterId := getClusterIDFromContext(ctx)

	clusteUser := f.fetcher.clusters[clusterId].Admin

	req := f.generateGetReq(clusteUser)

	json_data, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	cluster := f.fetcher.clusters[clusterId]

	getTokenURL := cluster.URL + ":5000" + "/v3/auth/tokens"

	resp, err := f.fetcher.client.Post(getTokenURL, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, handleErrorResponse(resp)
	}

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

func (f *TokenFetcher) Refresh(ctx context.Context, oldToken *models.Token) (*models.Token, error) {
	req := f.generateRefreshReq(oldToken)

	json_data, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	getTokenURL := cluster.URL + ":5000" + "/v3/auth/tokens"

	resp, err := f.fetcher.client.Post(getTokenURL, "application/json", bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, handleErrorResponse(resp)
	}

	newToken := &models.Token{}
	tokenRes := &GetTokenResponse{
		Token: newToken,
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenRes); err != nil {
		return nil, err
	}

	newToken.Value = resp.Header.Get("X-Subject-Token")

	return newToken, nil
}

func (f *TokenFetcher) generateRefreshReq(token *models.Token) *RefreshTokenRequest {
	methods := [1]string{"token"}

	req := &RefreshTokenRequest{
		Auth: &RefreshTokenAuth{
			Identity: &TokenIdentity{
				Methods: methods[:],
				Token: &Token{
					ID: token.Value,
				},
			},
		},
	}
	return req
}

func (f *TokenFetcher) generateGetReq(clusteUser *models.ClusterUser) *GetTokenRequest {
	methods := [1]string{"password"}

	req := &GetTokenRequest{
		Auth: &GetTokenAuth{
			Identity: &PasswordIdentity{
				Methods: methods[:],
				Password: &Password{
					User: &User{
						Name: clusteUser.Name,
						Domain: &Domain{
							ID: clusteUser.DomainID,
						},
						Password: clusteUser.Password,
					},
				},
			},
			Scope: &Scope{
				Project: &GetTokenProject{
					ID: clusteUser.ProjectID,
					Domain: &Domain{
						ID: clusteUser.DomainID,
					},
				},
			},
		},
	}

	return req
}
