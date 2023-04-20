package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RoleFetcher struct {
	fetcher *Fetcher
}

func (f *RoleFetcher) FetchByID(ctx context.Context, clusterId int, roleId string) (*models.Role, error) {
	cluster := f.fetcher.clusters[clusterId-1]

	fetchRoleURL := cluster.URL + ":5000" + "/v3/roles/" + roleId

	req, err := http.NewRequest("GET", fetchRoleURL, nil)
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

	role := &models.Role{}

	fetchRoleRes := &FetchRoleResponse{
		Role: role,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchRoleRes); err != nil {
		return nil, err
	}

	return fetchRoleRes.Role, nil
}

func (f *RoleFetcher) Create(ctx context.Context, clusterId int, role *models.Role) error {
	reqData := f.generateCreateReq(role)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	cluster := f.fetcher.clusters[clusterId-1]

	createRoleURL := cluster.URL + ":5000" + "/v3/roles"

	req, err := http.NewRequest("POST", createRoleURL, bytes.NewBuffer(json_data))
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
	log.Println(resp)

	createRoleRes := &CreateRoleResponse{
		Role: role,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createRoleRes); err != nil {
		return err
	}

	return nil

}

func (f *RoleFetcher) Delete(ctx context.Context, clusterId int, roleID string) error {
	cluster := f.fetcher.clusters[clusterId-1]

	deleteRoleURL := cluster.URL + ":5000" + "/v3/roles/" + roleID

	req, err := http.NewRequest("DELETE", deleteRoleURL, nil)
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

func (f *RoleFetcher) Update(ctx context.Context, clusterId int) {

}

func (f *RoleFetcher) getAdminToken(ctx context.Context, clusterId int) (*models.Token, error) {

	token, err := f.fetcher.Token().Get(ctx, clusterId, f.fetcher.clusters[clusterId-1].Admin)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (f *RoleFetcher) generateCreateReq(role *models.Role) *CreateRoleRequest {
	req := &CreateRoleRequest{
		Role: &Role{
			Name: role.Name,
		},
	}

	return req
}
