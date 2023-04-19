package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSecurityRuleFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	r := openstack.TestSecurityRule(t)

	err := fetcher.SecurityRule().Create(models.TestRequestContext(t), clusterID, r)
	assert.NoError(t, err)
	assert.NotEqual(t, r.ID, "")

	time.Sleep(1000)

	fetcher.SecurityRule().Delete(models.TestRequestContext(t), clusterID, r.ID)
}

func TestSecurityRuleFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	r := openstack.TestSecurityRule(t)

	fetcher.SecurityRule().Create(models.TestRequestContext(t), clusterID, r)

	time.Sleep(1000)

	err := fetcher.SecurityRule().Delete(models.TestRequestContext(t), clusterID, r.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, r.ID, "")
}

func TestSecurityRuleFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	r1 := openstack.TestSecurityRule(t)

	fetcher.SecurityRule().Create(models.TestRequestContext(t), clusterID, r1)

	time.Sleep(1000)

	r2, err := fetcher.SecurityRule().FetchByID(models.TestRequestContext(t), clusterID, r1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, r2.ID, "")

	fetcher.SecurityRule().Delete(models.TestRequestContext(t), clusterID, r2.ID)
}
