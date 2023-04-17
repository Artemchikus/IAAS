package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type UserFetcher struct {
	fetcher *Fetcher
}

func (f *UserFetcher) FetchByID(ctx context.Context, clusterId int, userId string) (*models.Account, error) {
	cluster := f.fetcher.clusters[clusterId-1]

	fetchUserURL := cluster.URL + "/v3/users/" + userId

	req, err := http.NewRequest("GET", fetchUserURL, nil)
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

	userResp := map[string]interface{}{}

	fetchUserResp := &FetchUserResponse{
		User: &userResp,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchUserResp); err != nil {
		return nil, err
	}

	user := &models.Account{
		OpenstackID: userResp["id"].(string),
		Name:        userResp["name"].(string),
		Email:       userResp["email"].(string),
		ProjectID:   userResp["default_project_id"].(string),
	}

	return user, nil
}

func (f *UserFetcher) Create(ctx context.Context, clusterId int, user *models.Account) error {
	description := "Project for user " + user.Email

	project := models.NewProject(user.Email, description)

	if err := f.fetcher.Project().Create(ctx, clusterId, project); err != nil {
		return err
	}

	reqDara := f.generateCreateReq(user, project.ID)

	json_data, err := json.Marshal(&reqDara)
	if err != nil {
		return err
	}

	cluster := f.fetcher.clusters[clusterId-1]

	createUserURL := cluster.URL + "/v3/users"

	req, err := http.NewRequest("POST", createUserURL, bytes.NewBuffer(json_data))
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

	userResp := map[string]interface{}{}

	createUserResp := &CreateUserResponse{
		User: &userResp,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createUserResp); err != nil {
		return err
	}

	user.OpenstackID = userResp["id"].(string)
	user.ProjectID = userResp["default_project_id"].(string)

	return nil
}

func (f *UserFetcher) Delete(ctx context.Context, clusterId int, userId string) error {
	cluster := f.fetcher.clusters[clusterId-1]

	fetchUserURL := cluster.URL + "/v3/users/" + userId

	req, err := http.NewRequest("DELETE", fetchUserURL, nil)
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

func (f *UserFetcher) Update(ctx context.Context, clusterId int) {}

func (f *UserFetcher) getAdminToken(ctx context.Context, clusterId int) (*models.Token, error) {

	token, err := f.fetcher.Token().Get(ctx, clusterId, f.fetcher.clusters[clusterId-1].Admin)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (f *UserFetcher) generateCreateReq(user *models.Account, projectID string) *CreateUserRequest {
	return &CreateUserRequest{
		User: &CreateUser{
			Name:      user.Name,
			DomainID:  "default",
			Password:  user.Password,
			Enabled:   true,
			ProjectID: projectID,
			Email:     user.Email,
		},
	}
}
