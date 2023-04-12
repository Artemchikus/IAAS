package postgres

import (
	"IAAS/internal/config"
	"IAAS/internal/models"
	"IAAS/internal/store"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Store struct {
	db                *sql.DB
	accountRepository *AccountRepository
	secretRepository  *SecretRepository
}

func New(db *sql.DB, config *config.ApiConfig) *Store {
	store := &Store{
		db: db,
	}
	if err := initialize(store, config); err != nil {
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

func initialize(store *Store, config *config.ApiConfig) error {
	if err := store.Account().Init(); err != nil {
		return err
	}
	log.Println("account initialized")

	jwtSecret := &models.Secret{
		Type:  "jwt",
		Value: config.JwtKey,
	}

	if err := store.Secret().Init(jwtSecret); err != nil {
		return err
	}
	log.Println("secret initialized")

	return nil
}
