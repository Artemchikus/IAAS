package postgres_test

import (
	"IAAS/internal/config"
	"IAAS/internal/models"
	"IAAS/internal/store"
	"IAAS/internal/store/postgres"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const databaseURL = "host=localhost port=5433 user=postgres password=iaas dbname=iaas-test sslmode=disable"

func TestAccountRepository_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret")

	config := config.NewConfig()
	config.JwtKey = "secretkey"

	initCtx := context.WithValue(context.Background(), models.CtxKeyRequestID, "test-initial-request")
	s := postgres.New(initCtx, db, config)

	ctx := context.WithValue(context.Background(), models.CtxKeyRequestID, "99999999-9999-9999-9999-999999999999")

	u := models.TestAccount(t)
	assert.NoError(t, s.Account().Create(ctx, u))
	assert.NotNil(t, u)
}

func TestAccountRepository_FindByEmail(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret")

	config := config.NewConfig()
	config.JwtKey = "secretkey"

	initCtx := context.WithValue(context.Background(), models.CtxKeyRequestID, "test-initial-request")
	s := postgres.New(initCtx, db, config)

	ctx := context.WithValue(context.Background(), models.CtxKeyRequestID, "99999999-9999-9999-9999-999999999999")

	u := models.TestAccount(t)
	_, err := s.Account().FindByEmail(ctx, u.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Account().Create(ctx, u)
	u, err = s.Account().FindByEmail(ctx, u.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestAccountRepository_FindByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret")

	config := config.NewConfig()
	config.JwtKey = "secretkey"

	initCtx := context.WithValue(context.Background(), models.CtxKeyRequestID, "test-initial-request")
	s := postgres.New(initCtx, db, config)

	ctx := context.WithValue(context.Background(), models.CtxKeyRequestID, "99999999-9999-9999-9999-999999999999")

	u1 := models.TestAccount(t)
	_, err := s.Account().FindByID(ctx, 1)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Account().Create(ctx, u1)
	u2, err := s.Account().FindByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
func TestAccountRepository_Update(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret")

	config := config.NewConfig()
	config.JwtKey = "secretkey"

	initCtx := context.WithValue(context.Background(), models.CtxKeyRequestID, "test-initial-request")
	s := postgres.New(initCtx, db, config)

	ctx := context.WithValue(context.Background(), models.CtxKeyRequestID, "99999999-9999-9999-9999-999999999999")

	u1 := models.TestAccount(t)
	_, err := s.Account().FindByID(ctx, 1)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Account().Create(ctx, u1)
	u2, err := s.Account().FindByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestAccountRepository_GetAll(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret")

	config := config.NewConfig()
	config.JwtKey = "secretkey"

	initCtx := context.WithValue(context.Background(), models.CtxKeyRequestID, "test-initial-request")
	s := postgres.New(initCtx, db, config)

	ctx := context.WithValue(context.Background(), models.CtxKeyRequestID, "99999999-9999-9999-9999-999999999999")

	_, err := s.Account().GetAll(ctx)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	u1 := models.TestAccount(t)
	s.Account().Create(ctx, u1)
	acs, err := s.Account().GetAll(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, acs)
}
