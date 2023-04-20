package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

type VolumeFetcher struct {
	fetcher *Fetcher
}

func (f *VolumeFetcher) FetchByID(ctx context.Context, clusterId int, volumeId string) (*models.Volume, error) {
	cluster := f.fetcher.clusters[clusterId-1]

	fetchVolumeURL := cluster.URL + ":8776" + "/v3/volumes/" + volumeId

	req, err := http.NewRequest("GET", fetchVolumeURL, nil)
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

	volumeResp := map[string]interface{}{}

	fetchVolumeRes := &FetchVolumeResponse{
		Volume: &volumeResp,
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
		ID:               volumeResp["id"].(string),
		Status:           volumeResp["status"].(string),
		Size:             int(volumeResp["size"].(float64)),
		AvailabilityZone: volumeResp["availability_zone"].(string),
		CreatedAt:        volumeResp["created_at"].(time.Time),
		UpdatedAt:        volumeResp["updated_at"].(time.Time),
		Type:             volumeResp["volume_type"].(string),
		AccountID:        volumeResp["user_id"].(string),
		Bootable:         volumeResp["bootable"].(bool),
		Encrypted:        volumeResp["encrypted"].(bool),
	}

	return volume, nil
}

func (f *VolumeFetcher) Create(ctx context.Context, clusterId int, volume *models.Volume) error {
	reqData := f.generateCreateReq(volume)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	cluster := f.fetcher.clusters[clusterId-1]

	createVolumeURL := cluster.URL + ":8776" + "/v3/volumes"

	req, err := http.NewRequest("POST", createVolumeURL, bytes.NewBuffer(json_data))
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

	volumeResp := map[string]interface{}{}

	createVolumeRes := &CreateVolumeResponse{
		Volume: &volumeResp,
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
	volume.Status = volumeResp["status"].(string)
	volume.Size = int(volumeResp["size"].(float64))
	volume.AvailabilityZone = volumeResp["availability_zone"].(string)
	volume.CreatedAt = volumeResp["created_at"].(time.Time)
	volume.UpdatedAt = volumeResp["created_at"].(time.Time)
	volume.Type = volumeResp["volume_type"].(string)
	volume.AccountID = volumeResp["user_id"].(string)
	volume.Bootable = volumeResp["bootable"].(bool)
	volume.Encrypted = volumeResp["encrypted"].(bool)

	return nil
}

func (f *VolumeFetcher) Delete(ctx context.Context, clusterId int, volumeID string) error {
	cluster := f.fetcher.clusters[clusterId-1]

	deleteVolumeURL := cluster.URL + ":8776" + "/v3/volumes/" + volumeID

	req, err := http.NewRequest("DELETE", deleteVolumeURL, nil)
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

	log.Println(resp)

	if resp.StatusCode != 202 {
		return errors.New("internal server error")
	}

	return nil
}

func (f *VolumeFetcher) Update(ctx context.Context, clusterId int) {

}

func (f *VolumeFetcher) getAdminToken(ctx context.Context, clusterId int) (*models.Token, error) {

	token, err := f.fetcher.Token().Get(ctx, clusterId, f.fetcher.clusters[clusterId-1].Admin)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (f *VolumeFetcher) generateCreateReq(volume *models.Volume) *CreateVolumeRequest {
	req := &CreateVolumeRequest{
		Volume: &Volume{
			Name:        volume.Name,
			Size:        volume.Size,
			Description: volume.Description,
			Type:        volume.Type,
		},
	}

	return req
}
