package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSubnetFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	sub := openstack.TestSubnet(t)
	n := openstack.TestNetwork(t)

	fetcher.Network().Create(models.TestRequestContext(t), clusterID, n)
	sub.NetworkID = n.ID

	err := fetcher.Subnet().Create(models.TestRequestContext(t), clusterID, sub)
	assert.NoError(t, err)
	assert.NotEqual(t, sub.ID, "")

	time.Sleep(1000)

	fetcher.Network().Delete(models.TestRequestContext(t), clusterID, n.ID)
}

func TestSubnetFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	sub := openstack.TestSubnet(t)
	n := openstack.TestNetwork(t)

	fetcher.Network().Create(models.TestRequestContext(t), clusterID, n)
	sub.NetworkID = n.ID

	fetcher.Subnet().Create(models.TestRequestContext(t), clusterID, sub)

	time.Sleep(1000)

	err := fetcher.Subnet().Delete(models.TestRequestContext(t), clusterID, sub.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, sub.ID, "")
	fetcher.Network().Delete(models.TestRequestContext(t), clusterID, n.ID)
}

func TestSubnetFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	sub1 := openstack.TestSubnet(t)
	n := openstack.TestNetwork(t)

	fetcher.Network().Create(models.TestRequestContext(t), clusterID, n)
	sub1.NetworkID = n.ID

	fetcher.Subnet().Create(models.TestRequestContext(t), clusterID, sub1)

	time.Sleep(1000)

	sub2, err := fetcher.Subnet().FetchByID(models.TestRequestContext(t), clusterID, sub1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, sub2.ID, "")

	fetcher.Network().Delete(models.TestRequestContext(t), clusterID, n.ID)
}
