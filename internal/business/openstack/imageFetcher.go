package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type ImageFetcher struct {
	fetcher *Fetcher
}

func (f *ImageFetcher) FetchByID(ctx context.Context, imageId string) (*models.Image, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchImageURL := cluster.URL + ":9292" + "/v2/images/" + imageId

	req, err := http.NewRequest("GET", fetchImageURL, nil)
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

	image := &models.Image{}

	if err := json.NewDecoder(resp.Body).Decode(&image); err != nil {
		return nil, err
	}

	return image, nil
}

func (f *ImageFetcher) Create(ctx context.Context, image *models.Image) error {
	reqData := f.generateCreateReq(image)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createImageURL := cluster.URL + ":9292" + "/v2/images"

	req, err := http.NewRequest("POST", createImageURL, bytes.NewBuffer(json_data))
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

	if err := json.NewDecoder(resp.Body).Decode(&image); err != nil {
		return err
	}

	return nil

}

func (f *ImageFetcher) Delete(ctx context.Context, imageId string) error {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	deleteImageURL := cluster.URL + ":9292" + "/v2/images/" + imageId

	req, err := http.NewRequest("DELETE", deleteImageURL, nil)
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

func (f *ImageFetcher) Upload(ctx context.Context, fileData []byte, imageId string) error {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	putImageDataURL := cluster.URL + ":9292" + "/v2/images/" + imageId + "/file"

	req, err := http.NewRequest("PUT", putImageDataURL, bytes.NewReader(fileData))
	if err != nil {
		return err
	}

	token := getTokenFromContext(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", token.Value)
	req.Header.Set("Content-Type", "application/octet-stream")

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

func (f *ImageFetcher) FetchAll(ctx context.Context) ([]*models.Image, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchImageURL := cluster.URL + ":9292" + "/v2/images"

	req, err := http.NewRequest("GET", fetchImageURL, nil)
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

	images := []*models.Image{}

	fetchImagesRes := &FetchImagesResponse{
		Images: &images,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchImagesRes); err != nil {
		return nil, err
	}

	return images, nil
}

func (f *ImageFetcher) generateCreateReq(image *models.Image) *CretaeImageRequest {
	req := &CretaeImageRequest{
		Name:            image.Name,
		DiskFormat:      image.DiskFormat,
		ContainerFormat: image.ContainerFormat,
		Visibility:      image.Visibility,
	}

	return req
}
