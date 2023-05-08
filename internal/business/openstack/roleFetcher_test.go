package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRoleFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	r := openstack.TestRole(t)

	err := fetcher.Role().Create(openstack.TestRequestContext(t, fetcher, clusterID), r)
	assert.NoError(t, err)
	assert.NotEqual(t, r.ID, "")

	time.Sleep(1000)

	fetcher.Role().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
}

func TestRoleFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	r := openstack.TestRole(t)

	fetcher.Role().Create(openstack.TestRequestContext(t, fetcher, clusterID), r)

	time.Sleep(1000)

	err := fetcher.Role().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, r.ID, "")
}

func TestRoleFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	r1 := openstack.TestRole(t)

	fetcher.Role().Create(openstack.TestRequestContext(t, fetcher, clusterID), r1)

	time.Sleep(1000)

	r2, err := fetcher.Role().FetchByID(openstack.TestRequestContext(t, fetcher, clusterID), r1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, r2.ID, "")

	fetcher.Role().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r2.ID)
}

func TestRoleFetcher_FetchByName(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	r1 := openstack.TestRole(t)

	fetcher.Role().Create(openstack.TestRequestContext(t, fetcher, clusterID), r1)

	time.Sleep(1000)

	r2, err := fetcher.Role().FetchByName(openstack.TestRequestContext(t, fetcher, clusterID), r1.Name)
	assert.NoError(t, err)
	assert.NotEqual(t, r2.ID, "")

	fetcher.Role().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r2.ID)
}

func TestRoleFetcher_FetchAll(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	r := openstack.TestRole(t)

	fetcher.Role().Create(openstack.TestRequestContext(t, fetcher, clusterID), r)

	time.Sleep(1000)

	rs, err := fetcher.Role().FetchAll(openstack.TestRequestContext(t, fetcher, clusterID))
	assert.NoError(t, err)
	assert.NotEmpty(t, rs)

	fetcher.Role().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
}
