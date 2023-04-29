package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFloatingIpFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	ip := &models.FloatingIp{}
	sub := openstack.TestPublicSubnet(t)
	n := openstack.TestPublicNetwork(t)

	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), n)
	sub.NetworkID = n.ID
	ip.NetworkID = n.ID

	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), sub)

	err := fetcher.FloatingIp().Create(openstack.TestRequestContext(t, fetcher, clusterID), ip)
	assert.NoError(t, err)
	assert.NotEqual(t, ip.ID, "")

	time.Sleep(1000)

	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), n.ID)
}

func TestFloatingIpFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	ip := &models.FloatingIp{}
	sub := openstack.TestPublicSubnet(t)
	n := openstack.TestPublicNetwork(t)

	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), n)
	sub.NetworkID = n.ID
	ip.NetworkID = n.ID

	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), sub)

	fetcher.FloatingIp().Create(openstack.TestRequestContext(t, fetcher, clusterID), ip)

	time.Sleep(1000)

	err := fetcher.FloatingIp().Delete(openstack.TestRequestContext(t, fetcher, clusterID), ip.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, ip.ID, "")
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), n.ID)
}

func TestFloatingIpFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	ip1 := &models.FloatingIp{}
	sub := openstack.TestPublicSubnet(t)
	n := openstack.TestPublicNetwork(t)

	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), n)
	sub.NetworkID = n.ID
	ip1.NetworkID = n.ID

	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), sub)

	fetcher.FloatingIp().Create(openstack.TestRequestContext(t, fetcher, clusterID), ip1)

	time.Sleep(1000)

	ip2, err := fetcher.FloatingIp().FetchByID(openstack.TestRequestContext(t, fetcher, clusterID), ip1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, ip2.ID, "")

	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), n.ID)
}
