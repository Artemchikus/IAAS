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

func (f *UserFetcher) FetchByID(ctx context.Context, userId string) (*models.Account, error) {
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

func (f *UserFetcher) Create(ctx context.Context, user *models.Account) error {
	description := "Project for user " + user.Email

	project := models.NewProject(user.Email, description)

	if err := f.fetcher.Project().Create(ctx, project); err != nil {
		return err
	}

	reqData := f.generateCreateReq(user, project.ID)

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
		return errors.New("internal server error")
	}

	return nil
}

func (f *UserFetcher) Update(ctx context.Context, userId string) {}

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
