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
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	ip := &models.FloatingIp{}
	sub := openstack.TestSubnet(t)
	n := openstack.TestNetwork(t)

	fetcher.Network().Create(models.TestRequestContext(t), clusterID, n)
	sub.NetworkID = n.ID
	ip.NetworkID = n.ID

	fetcher.Subnet().Create(models.TestRequestContext(t), clusterID, sub)

	err := fetcher.FloatingIp().Create(models.TestRequestContext(t), clusterID, ip)
	assert.NoError(t, err)
	assert.NotEqual(t, ip.ID, "")

	time.Sleep(1000)

	fetcher.Network().Delete(models.TestRequestContext(t), clusterID, n.ID)
}

func TestFloatingIpFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	ip := &models.FloatingIp{}
	sub := openstack.TestSubnet(t)
	n := openstack.TestNetwork(t)

	fetcher.Network().Create(models.TestRequestContext(t), clusterID, n)
	sub.NetworkID = n.ID
	ip.NetworkID = n.ID

	fetcher.Subnet().Create(models.TestRequestContext(t), clusterID, sub)

	fetcher.FloatingIp().Create(models.TestRequestContext(t), clusterID, ip)

	time.Sleep(1000)

	err := fetcher.FloatingIp().Delete(models.TestRequestContext(t), clusterID, ip.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, ip.ID, "")
	fetcher.Network().Delete(models.TestRequestContext(t), clusterID, n.ID)
}

func TestFloatingIpFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	ip1 := &models.FloatingIp{}
	sub := openstack.TestSubnet(t)
	n := openstack.TestNetwork(t)

	fetcher.Network().Create(models.TestRequestContext(t), clusterID, n)
	sub.NetworkID = n.ID
	ip1.NetworkID = n.ID

	fetcher.Subnet().Create(models.TestRequestContext(t), clusterID, sub)

	fetcher.FloatingIp().Create(models.TestRequestContext(t), clusterID, ip1)

	time.Sleep(1000)

	ip2, err := fetcher.FloatingIp().FetchByID(models.TestRequestContext(t), clusterID, ip1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, ip2.ID, "")

	fetcher.Network().Delete(models.TestRequestContext(t), clusterID, n.ID)
}
