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

type ImageFetcher struct {
	fetcher *Fetcher
}

func (f *ImageFetcher) FetchByID(ctx context.Context, clusterId int, imageId string) (*models.Image, error) {
	cluster := f.fetcher.clusters[clusterId-1]

	fetchImageURL := cluster.URL + ":9292" + "/v2/images/" + imageId

	req, err := http.NewRequest("GET", fetchImageURL, nil)
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

	image := &models.Image{}

	if err := json.NewDecoder(resp.Body).Decode(&image); err != nil {
		return nil, err
	}

	log.Println(image)

	return image, nil
}

func (f *ImageFetcher) Create(ctx context.Context, clusterId int, image *models.Image) error {
	reqDara := f.generateCreateReq(image)

	json_data, err := json.Marshal(&reqDara)
	if err != nil {
		return err
	}

	cluster := f.fetcher.clusters[clusterId-1]

	createImageURL := cluster.URL + ":9292" + "/v2/images"

	req, err := http.NewRequest("POST", createImageURL, bytes.NewBuffer(json_data))
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

	if err := json.NewDecoder(resp.Body).Decode(&image); err != nil {
		return err
	}

	if err := f.uploadImage(ctx, image.FileData, clusterId, image.ID); err != nil {
		return err
	}

	image.FileData = nil

	return nil

}

func (f *ImageFetcher) Delete(ctx context.Context, clusterId int, imageId string) error {
	cluster := f.fetcher.clusters[clusterId-1]

	deleteImageURL := cluster.URL + ":9292" + "/v2/images/" + imageId

	req, err := http.NewRequest("DELETE", deleteImageURL, nil)
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

func (f *ImageFetcher) Update(ctx context.Context, clusterId int) {

}

func (f *ImageFetcher) getAdminToken(ctx context.Context, clusterId int) (*models.Token, error) {

	token, err := f.fetcher.Token().Get(ctx, clusterId, f.fetcher.clusters[clusterId-1].Admin)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (f *ImageFetcher) uploadImage(ctx context.Context, fileData []byte, clusterId int, imageId string) error {
	cluster := f.fetcher.clusters[clusterId-1]

	putImageDataURL := cluster.URL + ":9292" + "/v2/images/" + imageId + "/file"

	req, err := http.NewRequest("PUT", putImageDataURL, bytes.NewReader(fileData))
	if err != nil {
		return err
	}

	token, err := f.getAdminToken(ctx, clusterId)
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

func (f *ImageFetcher) generateCreateReq(image *models.Image) *CretaeImageRequest {
	req := &CretaeImageRequest{
		Name:            image.Name,
		DiskFormat:      image.DiskFormat,
		ContainerFormat: image.ContainerFormat,
		Visibility:      image.Visibility,
	}

	return req
}
