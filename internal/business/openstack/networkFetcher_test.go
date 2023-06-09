package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNetworkFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	n := openstack.TestPrivateNetwork(t)

	err := fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), n)
	assert.NoError(t, err)
	assert.NotEqual(t, n.ID, "")

	time.Sleep(1000)

	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), n.ID)
}

func TestNetworkFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	n := openstack.TestPrivateNetwork(t)

	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), n)

	time.Sleep(1000)

	err := fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), n.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, n.ID, "")
}

func TestNetworkFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	n1 := openstack.TestPrivateNetwork(t)

	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), n1)

	time.Sleep(1000)

	n2, err := fetcher.Network().FetchByID(openstack.TestRequestContext(t, fetcher, clusterID), n1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, n2.ID, "")

	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), n2.ID)
}

func TestNetworkFetcher_FetchAll(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	n := openstack.TestPrivateNetwork(t)

	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), n)

	time.Sleep(1000)

	ns, err := fetcher.Network().FetchAll(openstack.TestRequestContext(t, fetcher, clusterID))
	assert.NoError(t, err)
	assert.NotEmpty(t, ns)

	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), n.ID)
}
