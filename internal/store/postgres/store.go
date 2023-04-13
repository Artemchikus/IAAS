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
	db                *sql.DB
	accountRepository *AccountRepository
	secretRepository  *SecretRepository
	logger            *zap.SugaredLogger
}

func New(ctx context.Context, db *sql.DB, config *config.ApiConfig) *Store {
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

func (s *Store) initialize(ctx context.Context, store *Store, config *config.ApiConfig) error {
	if err := store.Account().Init(ctx, config.Admin); err != nil {
		return err
	}

	sugar := s.logger.With(
		"request_id", ctx.Value(models.CtxKeyRequestID),
	)

	sugar.Infof("table account is initialized")

	jwtSecret := &models.Secret{
		Type:  "jwt",
		Value: config.JwtKey,
	}

	if err := store.Secret().Init(ctx, jwtSecret); err != nil {
		return err
	}
	sugar.Infof("table secret is initialized")

	return nil
}
