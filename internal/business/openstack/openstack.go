package openstack

import (
	"IAAS/internal/business"
	"IAAS/internal/config"
	"IAAS/internal/models"
	"context"
	"net/http"

	"go.uber.org/zap"
)

type Fetcher struct {
	logger        *zap.SugaredLogger
	client        *http.Client
	serverFetcher *ServerFetcher
	userFetcher   *UserFetcher
	tokenFetcher  *TokenFetcher
	clusters      []*models.Cluster
}

func New(ctx context.Context, config *config.ApiConfig) *Fetcher {
	log, _ := zap.NewProduction()
	defer log.Sync()
	sugar := log.Sugar()

	client := &http.Client{}

	return &Fetcher{
		logger:   sugar,
		client:   client,
		clusters: config.Clusters,
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
