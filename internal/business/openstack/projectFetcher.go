package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type ProjectFetcher struct {
	fetcher *Fetcher
}

func (f *ProjectFetcher) FetchByID(ctx context.Context, projectId string) (*models.Project, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchProjectURL := cluster.URL + ":5000" + "/v3/projects/" + projectId

	req, err := http.NewRequest("GET", fetchProjectURL, nil)
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

	project := &models.Project{}

	fetchProjectRes := &FetchProjectResponse{
		Project: project,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchProjectRes); err != nil {
		return nil, err
	}

	return fetchProjectRes.Project, nil
}

func (f *ProjectFetcher) Create(ctx context.Context, project *models.Project) error {
	reqData := f.generateCreateReq(project)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createProjectURL := cluster.URL + ":5000" + "/v3/projects"

	req, err := http.NewRequest("POST", createProjectURL, bytes.NewBuffer(json_data))
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

	createProjectRes := &CreateProjectResponse{
		Project: project,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createProjectRes); err != nil {
		return err
	}

	return nil

}

func (f *ProjectFetcher) Delete(ctx context.Context, projectID string) error {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	deleteProjectURL := cluster.URL + ":5000" + "/v3/projects/" + projectID

	req, err := http.NewRequest("DELETE", deleteProjectURL, nil)
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

func (f *ProjectFetcher) generateCreateReq(project *models.Project) *CreateProjectRequest {
	req := &CreateProjectRequest{
		Project: &CreateProject{
			DomainID:    project.DomainID,
			Name:        project.Name,
			Description: project.Description,
		},
	}

	return req
}
