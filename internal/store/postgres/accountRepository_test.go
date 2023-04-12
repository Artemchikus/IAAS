package postgres_test

import (
	"IAAS/internal/config"
	"IAAS/internal/models"
	"IAAS/internal/store"
	"IAAS/internal/store/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

const databaseURL = "host=localhost port=5433 user=postgres password=iaas dbname=iaas-test sslmode=disable"

func TestAccountRepository_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret")

	config := config.NewConfig()
	config.JwtKey = "secretkey"

	s := postgres.New(db, config)
	u := models.TestAccount(t)
	assert.NoError(t, s.Account().Create(u))
	assert.NotNil(t, u)
}

func TestAccountRepository_FindByEmail(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret")

	config := config.NewConfig()
	config.JwtKey = "secretkey"

	s := postgres.New(db, config)
	u := models.TestAccount(t)
	_, err := s.Account().FindByEmail(u.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Account().Create(u)
	u, err = s.Account().FindByEmail(u.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestAccountRepository_FindByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret")

	config := config.NewConfig()
	config.JwtKey = "secretkey"

	s := postgres.New(db, config)
	u1 := models.TestAccount(t)
	_, err := s.Account().FindByID(1)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Account().Create(u1)
	u2, err := s.Account().FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
func TestAccountRepository_Update(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret")

	config := config.NewConfig()
	config.JwtKey = "secretkey"

	s := postgres.New(db, config)
	u1 := models.TestAccount(t)
	_, err := s.Account().FindByID(1)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Account().Create(u1)
	u2, err := s.Account().FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestAccountRepository_GetAll(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret")

	config := config.NewConfig()
	config.JwtKey = "secretkey"

	s := postgres.New(db, config)
	_, err := s.Account().GetAll()
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u1 := models.TestAccount(t)
	s.Account().Create(u1)
	acs, err := s.Account().GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, acs)
}
