package postgres_test

import (
	"IAAS/internal/models"
	"IAAS/internal/store"
	"IAAS/internal/store/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClusterRepository_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := postgres.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	c := postgres.TestCluster(t)
	assert.NoError(t, s.Cluster().Create(models.TestRequestContext(t), c))
	assert.NotNil(t, c)
}

func TestClusterRepository_FindByLocation(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := postgres.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	c1 := postgres.TestCluster(t)
	_, err := s.Cluster().FindByLocation(models.TestRequestContext(t), c1.Location)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Cluster().Create(models.TestRequestContext(t), c1)
	c2, err := s.Cluster().FindByLocation(models.TestRequestContext(t), c1.Location)
	assert.NoError(t, err)
	assert.NotNil(t, c2)
}

func TestClusterRepository_FindByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := postgres.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	c1 := postgres.TestCluster(t)
	_, err := s.Cluster().FindByID(models.TestRequestContext(t), 2)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Cluster().Create(models.TestRequestContext(t), c1)
	c2, err := s.Cluster().FindByID(models.TestRequestContext(t), c1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, c2)
}
func TestClusterRepository_Update(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := postgres.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	c1 := postgres.TestCluster(t)
	_, err := s.Cluster().FindByID(models.TestRequestContext(t), 2)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.Cluster().Create(models.TestRequestContext(t), c1)
	c2, err := s.Cluster().FindByID(models.TestRequestContext(t), c1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, c2)
}

func TestClusterRepository_GetAll(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := postgres.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	err := s.Cluster().Delete(models.TestRequestContext(t), 1)
	assert.NoError(t, err)

	_, err = s.Cluster().GetAll(models.TestRequestContext(t))
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	c1 := postgres.TestCluster(t)
	s.Cluster().Create(models.TestRequestContext(t), c1)
	cls, err := s.Cluster().GetAll(models.TestRequestContext(t))
	assert.NoError(t, err)
	assert.NotNil(t, cls)
}
