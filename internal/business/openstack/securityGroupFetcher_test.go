package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSecurityGroupFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	sg := openstack.TestSecurityGroup(t)

	err := fetcher.SecurityGroup().Create(openstack.TestRequestContext(t, fetcher, clusterID), sg)
	assert.NoError(t, err)
	assert.NotEqual(t, sg.ID, "")

	time.Sleep(1000)

	fetcher.SecurityGroup().Delete(openstack.TestRequestContext(t, fetcher, clusterID), sg.ID)
}

func TestSecurityGroupFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	sg := openstack.TestSecurityGroup(t)

	fetcher.SecurityGroup().Create(openstack.TestRequestContext(t, fetcher, clusterID), sg)

	time.Sleep(1000)

	err := fetcher.SecurityGroup().Delete(openstack.TestRequestContext(t, fetcher, clusterID), sg.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, sg.ID, "")
}

func TestSecurityGroupFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	sg1 := openstack.TestSecurityGroup(t)

	fetcher.SecurityGroup().Create(openstack.TestRequestContext(t, fetcher, clusterID), sg1)

	time.Sleep(1000)

	sg2, err := fetcher.SecurityGroup().FetchByID(openstack.TestRequestContext(t, fetcher, clusterID), sg1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, sg2.ID, "")

	fetcher.SecurityGroup().Delete(openstack.TestRequestContext(t, fetcher, clusterID), sg2.ID)
}
