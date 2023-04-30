package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type UserFetcher struct {
	fetcher *Fetcher
}

func (f *UserFetcher) FetchByID(ctx context.Context, userId string) (*models.ClusterUser, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchUserURL := cluster.URL + ":5000" + "/v3/users/" + userId

	req, err := http.NewRequest("GET", fetchUserURL, nil)
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

	user := &models.ClusterUser{}

	fetchUserResp := &FetchUserResponse{
		User: user,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchUserResp); err != nil {
		return nil, err
	}

	return user, nil
}

func (f *UserFetcher) Create(ctx context.Context, user *models.ClusterUser) error {
	reqData := f.generateCreateReq(user)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createUserURL := cluster.URL + ":5000" + "/v3/users"

	req, err := http.NewRequest("POST", createUserURL, bytes.NewBuffer(json_data))
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

	if resp.StatusCode != 201 {
		return handleErrorResponse(resp)
	}

	createUserResp := &CreateUserResponse{
		User: user,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createUserResp); err != nil {
		return err
	}

	user.ClusterID = clusterId

	return nil
}

func (f *UserFetcher) Delete(ctx context.Context, userId string) error {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchUserURL := cluster.URL + ":5000" + "/v3/users/" + userId

	req, err := http.NewRequest("DELETE", fetchUserURL, nil)
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
		return handleErrorResponse(resp)
	}

	return nil
}

func (f *UserFetcher) FetchAll(ctx context.Context) ([]*models.ClusterUser, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchUserURL := cluster.URL + ":5000" + "/v3/users"

	req, err := http.NewRequest("GET", fetchUserURL, nil)
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

	users := []*models.ClusterUser{}

	fetchUsersRes := &FetchUsersResponse{
		Users: &users,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchUsersRes); err != nil {
		return nil, err
	}

	return users, nil
}

func (f *UserFetcher) generateCreateReq(user *models.ClusterUser) *CreateUserRequest {
	return &CreateUserRequest{
		User: &CreateUser{
			Name:        user.Name,
			DomainID:    user.DomainID,
			Password:    user.Password,
			ProjectID:   user.ProjectID,
			Email:       user.Email,
			Description: user.Description,
		},
	}
}
