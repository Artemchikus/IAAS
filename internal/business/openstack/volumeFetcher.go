package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type VolumeFetcher struct {
	fetcher *Fetcher
}

func (f *VolumeFetcher) FetchByID(ctx context.Context, volumeId string) (*models.Volume, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchVolumeURL := cluster.URL + ":8776" + "/v3/volumes/" + volumeId

	req, err := http.NewRequest("GET", fetchVolumeURL, nil)
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

	volumeResp := map[string]interface{}{}

	fetchVolumeRes := &FetchVolumeResponse{
		Volume: volumeResp,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchVolumeRes); err != nil {
		return nil, err
	}

	volumeResp["created_at"], err = time.Parse("2006-01-02T15:04:05.000000", volumeResp["created_at"].(string))
	if err != nil {
		return nil, err
	}

	volumeResp["updated_at"], err = time.Parse("2006-01-02T15:04:05.000000", volumeResp["updated_at"].(string))
	if err != nil {
		return nil, err
	}

	volumeResp["bootable"], err = strconv.ParseBool(volumeResp["bootable"].(string))
	if err != nil {
		return nil, err
	}

	volume := &models.Volume{
		ID:        volumeResp["id"].(string),
		Status:    volumeResp["status"].(string),
		Size:      int(volumeResp["size"].(float64)),
		CreatedAt: volumeResp["created_at"].(time.Time),
		UpdatedAt: volumeResp["updated_at"].(time.Time),
		TypeID:    volumeResp["volume_type"].(string),
		AccountID: volumeResp["user_id"].(string),
		Bootable:  volumeResp["bootable"].(bool),
		Name:      volumeResp["name"].(string),
	}

	return volume, nil
}

func (f *VolumeFetcher) Create(ctx context.Context, volume *models.Volume) error {
	reqData := f.generateCreateReq(volume)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	createVolumeURL := cluster.URL + ":8776" + "/v3/volumes"

	req, err := http.NewRequest("POST", createVolumeURL, bytes.NewBuffer(json_data))
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

	if resp.StatusCode != 202 {
		return handleErrorResponse(resp)
	}

	volumeResp := map[string]interface{}{}

	createVolumeRes := &CreateVolumeResponse{
		Volume: volumeResp,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createVolumeRes); err != nil {
		return err
	}

	volumeResp["created_at"], err = time.Parse("2006-01-02T15:04:05.000000", volumeResp["created_at"].(string))
	if err != nil {
		return err
	}

	volumeResp["bootable"], err = strconv.ParseBool(volumeResp["bootable"].(string))
	if err != nil {
		return err
	}

	volume.ID = volumeResp["id"].(string)
	volume.Name = volumeResp["name"].(string)
	volume.Status = volumeResp["status"].(string)
	volume.Size = int(volumeResp["size"].(float64))
	volume.CreatedAt = volumeResp["created_at"].(time.Time)
	volume.UpdatedAt = volumeResp["created_at"].(time.Time)
	volume.TypeID = volumeResp["volume_type"].(string)
	volume.AccountID = volumeResp["user_id"].(string)
	volume.Bootable = volumeResp["bootable"].(bool)

	return nil
}

func (f *VolumeFetcher) Delete(ctx context.Context, volumeID string) error {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	deleteVolumeURL := cluster.URL + ":8776" + "/v3/volumes/" + volumeID

	req, err := http.NewRequest("DELETE", deleteVolumeURL, nil)
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

	if resp.StatusCode != 202 {
		return handleErrorResponse(resp)
	}

	return nil
}

func (f *VolumeFetcher) FetchAll(ctx context.Context) ([]*models.Volume, error) {
	clusterId := getClusterIDFromContext(ctx)

	cluster := f.fetcher.clusters[clusterId]

	fetchVolumeURL := cluster.URL + ":8776" + "/v3/volumes"

	req, err := http.NewRequest("GET", fetchVolumeURL, nil)
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

	volumesResp := []map[string]interface{}{}

	fetchVolumesRes := &FetchVolumesResponse{
		Volumes: &volumesResp,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchVolumesRes); err != nil {
		return nil, err
	}

	volumes := []*models.Volume{}

	for _, v := range volumesResp {
		volume := &models.Volume{
			ID:   v["id"].(string),
			Name: v["name"].(string),
		}

		volumes = append(volumes, volume)
	}

	return volumes, nil
}

func (f *VolumeFetcher) generateCreateReq(volume *models.Volume) *CreateVolumeRequest {
	req := &CreateVolumeRequest{
		Volume: &Volume{
			Name:        volume.Name,
			Size:        volume.Size,
			Description: volume.Description,
			TypeID:      volume.TypeID,
			Bootable:    volume.Bootable,
		},
	}

	return req
}
