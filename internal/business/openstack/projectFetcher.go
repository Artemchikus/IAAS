package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type ProjectFetcher struct {
	fetcher *Fetcher
}

func (f *ProjectFetcher) FetchByID(ctx context.Context, clusterId int, projectId string) (*models.Project, error) {
	cluster := f.fetcher.clusters[clusterId-1]

	fetchProjectURL := cluster.URL + "/v3/projects/" + projectId

	req, err := http.NewRequest("GET", fetchProjectURL, nil)
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

	project := &models.Project{}

	fetchProjectRes := &FetchProjectResponse{
		Project: project,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchProjectRes); err != nil {
		return nil, err
	}

	return fetchProjectRes.Project, nil
}

func (f *ProjectFetcher) Create(ctx context.Context, clusterId int, project *models.Project) error {
	reqDara := f.generateCreateReq(project)

	json_data, err := json.Marshal(&reqDara)
	if err != nil {
		return err
	}

	cluster := f.fetcher.clusters[clusterId-1]

	createProjectURL := cluster.URL + "/v3/projects"

	req, err := http.NewRequest("POST", createProjectURL, bytes.NewBuffer(json_data))
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

	createProjectRes := &CreateProjectResponse{
		Project: project,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createProjectRes); err != nil {
		return err
	}

	return nil

}

func (f *ProjectFetcher) Delete(ctx context.Context, clusterId int, projectID string) error {
	cluster := f.fetcher.clusters[clusterId-1]

	deleteProjectURL := cluster.URL + "/v3/projects/" + projectID

	req, err := http.NewRequest("DELETE", deleteProjectURL, nil)
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

func (f *ProjectFetcher) Update(ctx context.Context, clusterId int) {

}

func (f *ProjectFetcher) getAdminToken(ctx context.Context, clusterId int) (*models.Token, error) {

	token, err := f.fetcher.Token().Get(ctx, clusterId, f.fetcher.clusters[clusterId-1].Admin)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (f *ProjectFetcher) generateCreateReq(project *models.Project) *CreateProjectRequest {
	req := &CreateProjectRequest{
		Project: &CreateProject{
			DomainID:    project.DomainID,
			Name:        project.Name,
			Description: project.Description,
			Enabled:     project.Enabled,
		},
	}

	return req
}
