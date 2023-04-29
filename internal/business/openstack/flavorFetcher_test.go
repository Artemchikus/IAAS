package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFlavorFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	f := openstack.TestFlavor(t)

	err := fetcher.Flavor().Create(openstack.TestRequestContext(t, fetcher, clusterID), f)
	assert.NoError(t, err)
	assert.NotEqual(t, f.ID, "")

	time.Sleep(1000)

	fetcher.Flavor().Delete(openstack.TestRequestContext(t, fetcher, clusterID), f.ID)
}

func TestFlavorFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	f := openstack.TestFlavor(t)

	fetcher.Flavor().Create(openstack.TestRequestContext(t, fetcher, clusterID), f)

	time.Sleep(1000)

	err := fetcher.Flavor().Delete(openstack.TestRequestContext(t, fetcher, clusterID), f.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, f.ID, "")
}

func TestFlavorFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	f1 := openstack.TestFlavor(t)

	fetcher.Flavor().Create(openstack.TestRequestContext(t, fetcher, clusterID), f1)

	time.Sleep(1000)

	f2, err := fetcher.Flavor().FetchByID(openstack.TestRequestContext(t, fetcher, clusterID), f1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, f2.ID, "")

	fetcher.Flavor().Delete(openstack.TestRequestContext(t, fetcher, clusterID), f2.ID)
}
