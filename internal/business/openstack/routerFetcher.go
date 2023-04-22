package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type RouterFetcher struct {
	fetcher *Fetcher
}

func (f *RouterFetcher) FetchByID(ctx context.Context, routerId string) (*models.Router, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchRouterURL := cluster.URL + ":9696" + "/v2.0/routers/" + routerId

	req, err := http.NewRequest("GET", fetchRouterURL, nil)
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

	router := &models.Router{}

	fetchRouterRes := &FetchRouterResponse{
		Router: router,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchRouterRes); err != nil {
		return nil, err
	}

	return fetchRouterRes.Router, nil
}

func (f *RouterFetcher) Create(ctx context.Context, router *models.Router) error {
	reqData := f.generateCreateReq(router)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createRouterURL := cluster.URL + ":9696" + "/v2.0/routers"

	req, err := http.NewRequest("POST", createRouterURL, bytes.NewBuffer(json_data))
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

	createRouterRes := &CreateRouterResponse{
		Router: router,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createRouterRes); err != nil {
		return err
	}

	return nil

}

func (f *RouterFetcher) Delete(ctx context.Context, routerID string) error {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	deleteRouterURL := cluster.URL + ":9696" + "/v2.0/routers/" + routerID

	req, err := http.NewRequest("DELETE", deleteRouterURL, nil)
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

func (f *RouterFetcher) Update(ctx context.Context, routerId string) {

}

func (f *RouterFetcher) generateCreateReq(router *models.Router) *CreateRouterRequest {
	req := &CreateRouterRequest{
		Router: &Router{
			Name:         router.Name,
			AdminStateUp: true,
		},
	}

	return req
}
