package postgres

import (
	"IAAS/internal/config"
	"IAAS/internal/models"
	"IAAS/internal/store"
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Store struct {
	db                    *sql.DB
	accountRepository     *AccountRepository
	secretRepository      *SecretRepository
	clusterRepository     *ClusterRepository
	clusterUserRepository *ClusterUserRepository
	logger                *zap.SugaredLogger
}

func NewStore(ctx context.Context, db *sql.DB, config *config.ApiConfig) *Store {
	zapLog, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer zapLog.Sync()
	sugar := zapLog.Sugar()

	store := &Store{
		db:     db,
		logger: sugar,
	}

	if err := store.initialize(ctx, store, config); err != nil {
		log.Fatal(err)
	}
	return store
}

func (s *Store) Account() store.AccountRepository {
	if s.accountRepository != nil {
		return s.accountRepository
	}

	s.accountRepository = &AccountRepository{
		store: s,
	}

	return s.accountRepository
}

func (s *Store) Secret() store.SecretRepository {
	if s.secretRepository != nil {
		return s.secretRepository
	}

	s.secretRepository = &SecretRepository{
		store: s,
	}

	return s.secretRepository
}

func (s *Store) Cluster() store.ClusterRepository {
	if s.clusterRepository != nil {
		return s.clusterRepository
	}

	s.clusterRepository = &ClusterRepository{
		store: s,
	}

	return s.clusterRepository
}

func (s *Store) ClusterUser() store.ClusterUserRepository {
	if s.clusterUserRepository != nil {
		return s.clusterUserRepository
	}

	s.clusterUserRepository = &ClusterUserRepository{
		store: s,
	}

	return s.clusterUserRepository
}

func (s *Store) initialize(ctx context.Context, store *Store, config *config.ApiConfig) error {
	if err := store.Account().Init(ctx, config.Admin); err != nil {
		return err
	}

	s.logger = s.logger.With(
		"request_id", ctx.Value(models.CtxKeyRequestID),
	)

	s.logger.Infof("table account is initialized")

	jwtSecret := &models.Secret{
		Type:  "jwt",
		Value: config.JwtKey,
	}

	if err := store.Secret().Init(ctx, jwtSecret); err != nil {
		return err
	}
	s.logger.Infof("table secret is initialized")

	if err := store.Cluster().Init(ctx, config.Clusters); err != nil {
		return err
	}
	s.logger.Infof("table cluster is initialized")

	if err := store.ClusterUser().Init(ctx, config.Clusters); err != nil {
		return err
	}
	s.logger.Infof("table clusterUser is initialized")

	return nil
}
