package openstack

import (
	"IAAS/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type KeyPairFetcher struct {
	fetcher *Fetcher
}

func (f *KeyPairFetcher) FetchByID(ctx context.Context, clusterId int, keyPairId string) (*models.KeyPair, error) {
	cluster := f.fetcher.clusters[clusterId-1]

	fetchKeyPairURL := cluster.URL + ":8774" + "/v2.1/os-keypairs/" + keyPairId

	req, err := http.NewRequest("GET", fetchKeyPairURL, nil)
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

	keyPairResp := map[string]interface{}{}

	fetchKeyPairRes := &FetchKeyPairResponse{
		KeyPair: &keyPairResp,
	}

	if err := json.NewDecoder(resp.Body).Decode(&fetchKeyPairRes); err != nil {
		return nil, err
	}

	keyPairResp["created_at"], err = time.Parse("2006-01-02T15:04:05.000000", keyPairResp["created_at"].(string))
	if err != nil {
		return nil, err
	}

	keyPair := &models.KeyPair{
		IsDeleted:   keyPairResp["deleted"].(bool),
		Fingerprint: keyPairResp["fingerprint"].(string),
		PublicKey:   keyPairResp["public_key"].(string),
		AccountID:   keyPairResp["user_id"].(string),
		ID:          keyPairResp["name"].(string),
		Type:        "ssh", // TODO add x509 support
		Name:        keyPairResp["name"].(string),
		CreatedAt:   keyPairResp["created_at"].(time.Time),
	}

	return keyPair, nil
}

func (f *KeyPairFetcher) Create(ctx context.Context, clusterId int, keyPair *models.KeyPair) error {
	reqData := f.generateCreateReq(keyPair)

	json_data, err := json.Marshal(&reqData)
	if err != nil {
		return err
	}

	cluster := f.fetcher.clusters[clusterId-1]

	createKeyPairURL := cluster.URL + ":8774" + "/v2.1/os-keypairs"

	req, err := http.NewRequest("POST", createKeyPairURL, bytes.NewBuffer(json_data))
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

	createKeyPairRes := &CreateKeyPairResponse{
		KeyPair: keyPair,
	}

	if err := json.NewDecoder(resp.Body).Decode(&createKeyPairRes); err != nil {
		return err
	}

	keyPair.ID = keyPair.Name

	return nil

}

func (f *KeyPairFetcher) Delete(ctx context.Context, clusterId int, keyPairId string) error {
	cluster := f.fetcher.clusters[clusterId-1]

	deleteKeyPairURL := cluster.URL + ":8774" + "/v2.1/os-keypairs/" + keyPairId

	req, err := http.NewRequest("DELETE", deleteKeyPairURL, nil)
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

	if resp.StatusCode != 202 {
		return errors.New("internal server error")
	}

	return nil
}

func (f *KeyPairFetcher) Update(ctx context.Context, clusterId int) {

}

func (f *KeyPairFetcher) getAdminToken(ctx context.Context, clusterId int) (*models.Token, error) {

	token, err := f.fetcher.Token().Get(ctx, clusterId, f.fetcher.clusters[clusterId-1].Admin)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (f *KeyPairFetcher) generateCreateReq(keyPair *models.KeyPair) *CreateKeyPairRequest {
	req := &CreateKeyPairRequest{
		KeyPair: &KeyPair{
			Name:      keyPair.Name,
			PublicKey: keyPair.PublicKey,
		},
	}

	return req
}
