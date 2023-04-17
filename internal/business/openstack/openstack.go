package openstack

import (
	"IAAS/internal/business"
	"IAAS/internal/config"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"context"
	"net/http"

	"go.uber.org/zap"
)

type Fetcher struct {
	logger         *zap.SugaredLogger
	client         *http.Client
	serverFetcher  *ServerFetcher
	userFetcher    *UserFetcher
	tokenFetcher   *TokenFetcher
	projectFetcher *ProjectFetcher
	imageFetcher   *ImageFetcher
	clusters       []*models.Cluster
}

func New(ctx context.Context, config *config.ApiConfig, store *postgres.Store) *Fetcher {
	log, _ := zap.NewProduction()
	defer log.Sync()
	sugar := log.Sugar()

	client := &http.Client{}

	clusters, err := store.Cluster().GetAll(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &Fetcher{
		logger:   sugar,
		client:   client,
		clusters: clusters,
	}
}

func (f *Fetcher) Server() business.ServerFetcher {
	if f.serverFetcher != nil {
		return f.serverFetcher
	}

	f.serverFetcher = &ServerFetcher{
		fetcher: f,
	}

	return f.serverFetcher
}

func (f *Fetcher) User() business.UserFetcher {
	if f.userFetcher != nil {
		return f.userFetcher
	}

	f.userFetcher = &UserFetcher{
		fetcher: f,
	}

	return f.userFetcher
}

func (f *Fetcher) Token() business.TokenFetcher {
	if f.tokenFetcher != nil {
		return f.tokenFetcher
	}

	f.tokenFetcher = &TokenFetcher{
		fetcher: f,
	}

	return f.tokenFetcher
}

func (f *Fetcher) Project() business.ProjectFetcher {
	if f.projectFetcher != nil {
		return f.projectFetcher
	}

	f.projectFetcher = &ProjectFetcher{
		fetcher: f,
	}

	return f.projectFetcher
}

func (f *Fetcher) Image() business.ImageFetcher {
	if f.imageFetcher != nil {
		return f.imageFetcher
	}

	f.imageFetcher = &ImageFetcher{
		fetcher: f,
	}

	return f.imageFetcher
}
