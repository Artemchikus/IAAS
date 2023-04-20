package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVolumeFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	v := openstack.TestVolume(t)

	err := fetcher.Volume().Create(models.TestRequestContext(t), clusterID, v)
	assert.NoError(t, err)
	assert.NotEqual(t, v.ID, "")

	time.Sleep(100000000)

	fetcher.Volume().Delete(models.TestRequestContext(t), clusterID, v.ID)
}

func TestVolumeFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	v := openstack.TestVolume(t)

	fetcher.Volume().Create(models.TestRequestContext(t), clusterID, v)

	time.Sleep(100000000)

	err := fetcher.Volume().Delete(models.TestRequestContext(t), clusterID, v.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, v.ID, "")
}

func TestVolumeFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	v1 := openstack.TestVolume(t)

	fetcher.Volume().Create(models.TestRequestContext(t), clusterID, v1)

	time.Sleep(100000000)

	v2, err := fetcher.Volume().FetchByID(models.TestRequestContext(t), clusterID, v1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, v2.ID, "")

	fetcher.Volume().Delete(models.TestRequestContext(t), clusterID, v2.ID)
}
