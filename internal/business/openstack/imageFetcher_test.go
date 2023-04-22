package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestImageFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	i := openstack.TestImage(t)

	err := fetcher.Image().Create(openstack.TestRequestContext(t, fetcher, clusterID), i)
	assert.NoError(t, err)
	assert.NotEqual(t, i.ID, "")
	assert.NotEqual(t, i.OwnerID, "")

	time.Sleep(1000)

	fetcher.Image().Delete(openstack.TestRequestContext(t, fetcher, clusterID), i.ID)
}

func TestImageFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	i := openstack.TestImage(t)

	fetcher.Image().Create(openstack.TestRequestContext(t, fetcher, clusterID), i)

	time.Sleep(1000)

	err := fetcher.Image().Delete(openstack.TestRequestContext(t, fetcher, clusterID), i.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, i.ID, "")
}

func TestImageFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	i1 := openstack.TestImage(t)

	fetcher.Image().Create(openstack.TestRequestContext(t, fetcher, clusterID), i1)

	time.Sleep(1000)

	i2, err := fetcher.Image().FetchByID(openstack.TestRequestContext(t, fetcher, clusterID), i1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, i2.ID, "")
	assert.NotEqual(t, i2.OwnerID, "")

	fetcher.Image().Delete(openstack.TestRequestContext(t, fetcher, clusterID), i2.ID)
}
