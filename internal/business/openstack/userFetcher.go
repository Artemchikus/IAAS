package openstack

import (
	"context"
	"log"
)

type UserFetcher struct {
	fetcher *Fetcher
}

func (f *UserFetcher) FetchByID(ctx context.Context, clusterId int) {}

func (f *UserFetcher) Create(ctx context.Context, clusterId int) {
	token, err := f.getAdminToken(ctx, clusterId)
	log.Println(token)
	if err != nil {
		log.Println(err)
	}
}

func (f *UserFetcher) Delete(ctx context.Context, clusterId int) {}

func (f *UserFetcher) Update(ctx context.Context, clusterId int) {}

func (f *UserFetcher) getAdminToken(ctx context.Context, clusterId int) (string, error) {

	token, err := f.fetcher.Token().Get(ctx, clusterId, f.fetcher.clusters[1].Admin)
	if err != nil {
		return "", err
	}

	return token.Value, nil
}
