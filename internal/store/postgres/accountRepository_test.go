package postgres_test

import (
	"IAAS/internal/models"
	"IAAS/internal/store"
	"IAAS/internal/store/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

const databaseURL = "host=localhost port=5433 user=postgres password=iaas dbname=iaas-test sslmode=disable"

func TestAccountRepository_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := postgres.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	u := models.TestAccount(t)
	assert.NoError(t, s.Account().Create(models.TestRequestContext(t), u))
	assert.NotNil(t, u)
}

func TestAccountRepository_FindByEmail(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := postgres.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	a1 := models.TestAccount(t)
	_, err := s.Account().FindByEmail(models.TestRequestContext(t), a1.Email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Account().Create(models.TestRequestContext(t), a1)
	a2, err := s.Account().FindByEmail(models.TestRequestContext(t), a1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, a2)
}

func TestAccountRepository_FindByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := postgres.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	a1 := models.TestAccount(t)
	_, err := s.Account().FindByID(models.TestRequestContext(t), 2)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Account().Create(models.TestRequestContext(t), a1)
	a2, err := s.Account().FindByID(models.TestRequestContext(t), a1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, a2)
}
func TestAccountRepository_Update(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := postgres.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	a1 := models.TestAccount(t)
	_, err := s.Account().FindByID(models.TestRequestContext(t), 2)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Account().Create(models.TestRequestContext(t), a1)
	a2, err := s.Account().FindByID(models.TestRequestContext(t), a1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, a2)
}

func TestAccountRepository_GetAll(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := postgres.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	err := s.Account().Delete(models.TestRequestContext(t), 1)
	assert.NoError(t, err)

	_, err = s.Account().GetAll(models.TestRequestContext(t))
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	a1 := models.TestAccount(t)
	s.Account().Create(models.TestRequestContext(t), a1)
	as, err := s.Account().GetAll(models.TestRequestContext(t))
	assert.NoError(t, err)
	assert.NotNil(t, as)
}
