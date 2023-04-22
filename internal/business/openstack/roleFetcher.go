package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type RoleFetcher struct {
	fetcher *Fetcher
}

func (f *RoleFetcher) FetchByID(ctx context.Context, roleId string) (*models.Role, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchRoleURL := cluster.URL + ":5000" + "/v3/roles/" + roleId

	req, err := http.NewRequest("GET", fetchRoleURL, nil)
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

	role := &models.Role{}

	fetchRoleRes := &FetchRoleResponse{
		Role: role,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchRoleRes); err != nil {
		return nil, err
	}

	return fetchRoleRes.Role, nil
}

func (f *RoleFetcher) Create(ctx context.Context, role *models.Role) error {
	reqData := f.generateCreateReq(role)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createRoleURL := cluster.URL + ":5000" + "/v3/roles"

	req, err := http.NewRequest("POST", createRoleURL, bytes.NewBuffer(json_data))
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

	createRoleRes := &CreateRoleResponse{
		Role: role,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createRoleRes); err != nil {
		return err
	}

	return nil

}

func (f *RoleFetcher) Delete(ctx context.Context, roleID string) error {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	deleteRoleURL := cluster.URL + ":5000" + "/v3/roles/" + roleID

	req, err := http.NewRequest("DELETE", deleteRoleURL, nil)
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

func (f *RoleFetcher) Update(ctx context.Context, roleId string) {

}

func (f *RoleFetcher) generateCreateReq(role *models.Role) *CreateRoleRequest {
	req := &CreateRoleRequest{
		Role: &Role{
			Name: role.Name,
		},
	}

	return req
}
