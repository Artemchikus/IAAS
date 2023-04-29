package postgres_test

import (
	"IAAS/internal/models"
	"IAAS/internal/store"
	"IAAS/internal/store/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretRepository_FindByType(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := postgres.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	_, err := s.Secret().FindByType(models.TestRequestContext(t), "test")
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	sec, err := s.Secret().FindByType(models.TestRequestContext(t), "jwt")
	assert.NoError(t, err)
	assert.NotNil(t, sec)
}
