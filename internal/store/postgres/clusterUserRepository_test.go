package postgres_test

import (
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClusterUserRepository_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := postgres.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	u := postgres.TestClusterUser(t)
	assert.NoError(t, s.ClusterUser().Create(models.TestRequestContext(t), u))
	assert.NotNil(t, u)
}

func TestClusterUserRepository_FindByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := postgres.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	a1 := postgres.TestClusterUser(t)

	s.ClusterUser().Create(models.TestRequestContext(t), a1)
	a2, err := s.ClusterUser().FindByID(models.TestRequestContext(t), a1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, a2)
}

func TestClusterUserRepository_FindByAccountAndClusterID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := postgres.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	a1 := postgres.TestClusterUser(t)

	s.ClusterUser().Create(models.TestRequestContext(t), a1)
	a2, err := s.ClusterUser().FindByEmailAndClusterID(models.TestRequestContext(t), a1.Email, a1.ClusterID)
	assert.NoError(t, err)
	assert.NotNil(t, a2)
}

func TestClusterUserRepository_Update(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := postgres.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	a1 := postgres.TestClusterUser(t)

	s.ClusterUser().Create(models.TestRequestContext(t), a1)
	a2, err := s.ClusterUser().FindByID(models.TestRequestContext(t), a1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, a2)
}

func TestClusterUserRepository_FindByAccountID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := postgres.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	a1 := postgres.TestClusterUser(t)
	s.ClusterUser().Create(models.TestRequestContext(t), a1)
	as, err := s.ClusterUser().FindByAccountID(models.TestRequestContext(t), a1.AccountID)
	assert.NoError(t, err)
	assert.NotNil(t, as)
}

func TestClusterUserRepository_FindByClusterID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := postgres.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	a1 := postgres.TestClusterUser(t)
	s.ClusterUser().Create(models.TestRequestContext(t), a1)
	as, err := s.ClusterUser().FindByClusterID(models.TestRequestContext(t), a1.ClusterID)
	assert.NoError(t, err)
	assert.NotNil(t, as)
}

func TestClusterUserRepository_GetAll(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := postgres.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	a1 := postgres.TestClusterUser(t)
	s.ClusterUser().Create(models.TestRequestContext(t), a1)
	as, err := s.ClusterUser().GetAll(models.TestRequestContext(t))
	assert.NoError(t, err)
	assert.NotNil(t, as)
}
