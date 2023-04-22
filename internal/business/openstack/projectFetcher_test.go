package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProjectFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	p := openstack.TestProject(t)

	err := fetcher.Project().Create(openstack.TestRequestContext(t, fetcher, clusterID), p)
	assert.NoError(t, err)
	assert.NotEqual(t, p.ID, "")

	time.Sleep(1000)

	fetcher.Project().Delete(openstack.TestRequestContext(t, fetcher, clusterID), p.ID)
}

func TestProjectFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	p := openstack.TestProject(t)

	fetcher.Project().Create(openstack.TestRequestContext(t, fetcher, clusterID), p)

	time.Sleep(1000)

	err := fetcher.Project().Delete(openstack.TestRequestContext(t, fetcher, clusterID), p.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, p.ID, "")
}

func TestProjectFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	p1 := openstack.TestProject(t)

	fetcher.Project().Create(openstack.TestRequestContext(t, fetcher, clusterID), p1)

	time.Sleep(1000)

	p2, err := fetcher.Project().FetchByID(openstack.TestRequestContext(t, fetcher, clusterID), p1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, p2.ID, "")

	fetcher.Project().Delete(openstack.TestRequestContext(t, fetcher, clusterID), p2.ID)
}
