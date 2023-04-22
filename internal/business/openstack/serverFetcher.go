package openstack

import (
	"IAAS/internal/models"
	"context"
)

type ServerFetcher struct {
	fetcher *Fetcher
}

func (f *ServerFetcher) FetchByID(ctx context.Context, serverId string) (*models.Server, error) {
	return nil, nil
}

func (f *ServerFetcher) Create(ctx context.Context, server *models.Server) error {
	return nil
}

func (f *ServerFetcher) Delete(ctx context.Context, serverId string) error {
	return nil
}

func (f *ServerFetcher) Update(ctx context.Context, serverId string) {}
