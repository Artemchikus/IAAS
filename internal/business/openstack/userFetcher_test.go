package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	p := openstack.TestProject(t)

	fetcher.Project().Create(openstack.TestRequestContext(t, fetcher, clusterID), p)

	u := openstack.TestClusterUser(t)
	u.ProjectID = p.ID

	err := fetcher.User().Create(openstack.TestRequestContext(t, fetcher, clusterID), u)
	assert.NoError(t, err)
	assert.NotEqual(t, u.ID, "")
	assert.NotEqual(t, u.ProjectID, "")

	time.Sleep(1000)

	fetcher.User().Delete(openstack.TestRequestContext(t, fetcher, clusterID), u.ID)
	fetcher.Project().Delete(openstack.TestRequestContext(t, fetcher, clusterID), u.ProjectID)
}

func TestUserFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	p := openstack.TestProject(t)

	fetcher.Project().Create(openstack.TestRequestContext(t, fetcher, clusterID), p)

	u1 := openstack.TestClusterUser(t)
	u1.ProjectID = p.ID

	fetcher.User().Create(openstack.TestRequestContext(t, fetcher, clusterID), u1)

	time.Sleep(1000)

	u2, err := fetcher.User().FetchByID(openstack.TestRequestContext(t, fetcher, clusterID), u1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, u2.ID, "")

	fetcher.User().Delete(openstack.TestRequestContext(t, fetcher, clusterID), u2.ID)
	fetcher.Project().Delete(openstack.TestRequestContext(t, fetcher, clusterID), u2.ProjectID)
}

func TestUserFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	p := openstack.TestProject(t)

	fetcher.Project().Create(openstack.TestRequestContext(t, fetcher, clusterID), p)

	u := openstack.TestClusterUser(t)
	u.ProjectID = p.ID

	err := fetcher.User().Create(openstack.TestRequestContext(t, fetcher, clusterID), u)

	time.Sleep(1000)

	fetcher.User().Delete(openstack.TestRequestContext(t, fetcher, clusterID), u.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, u.ID, "")
	fetcher.Project().Delete(openstack.TestRequestContext(t, fetcher, clusterID), u.ProjectID)
}
