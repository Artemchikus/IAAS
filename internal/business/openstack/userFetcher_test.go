package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	u := models.TestAccount(t)

	err := fetcher.User().Create(models.TestRequestContext(t), clusterID, u)
	assert.NoError(t, err)
	assert.NotEqual(t, u.OpenstackID, "")
	assert.NotEqual(t, u.ProjectID, "")

	log.Println(u)

	time.Sleep(1000)

	fetcher.User().Delete(models.TestRequestContext(t), clusterID, u.OpenstackID)
	fetcher.Project().Delete(models.TestRequestContext(t), clusterID, u.ProjectID)
}

func TestUserFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	u1 := models.TestAccount(t)

	fetcher.User().Create(models.TestRequestContext(t), clusterID, u1)

	time.Sleep(1000)

	u2, err := fetcher.User().FetchByID(models.TestRequestContext(t), clusterID, u1.OpenstackID)
	assert.NoError(t, err)
	assert.NotEqual(t, u2.ID, "")

	fetcher.User().Delete(models.TestRequestContext(t), clusterID, u2.OpenstackID)
	fetcher.Project().Delete(models.TestRequestContext(t), clusterID, u2.ProjectID)
}

func TestUserFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	u := models.TestAccount(t)

	err := fetcher.User().Create(models.TestRequestContext(t), clusterID, u)

	time.Sleep(1000)

	fetcher.User().Delete(models.TestRequestContext(t), clusterID, u.OpenstackID)
	assert.NoError(t, err)
	assert.NotEqual(t, u.OpenstackID, "")
	fetcher.Project().Delete(models.TestRequestContext(t), clusterID, u.ProjectID)
}
