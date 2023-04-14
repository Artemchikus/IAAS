package openstack

import "context"

type ServerFetcher struct {
	fetcher *Fetcher
}

func (f *ServerFetcher) FetchByID(ctx context.Context, clusterId int) {}

func (f *ServerFetcher) Create(ctx context.Context, clusterId int) {}

func (f *ServerFetcher) Delete(ctx context.Context, clusterId int) {}

func (f *ServerFetcher) Update(ctx context.Context, clusterId int) {}
