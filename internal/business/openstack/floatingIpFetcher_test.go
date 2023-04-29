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

func TestFloatinIpFetcher_AddToPort(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	prn := openstack.TestPrivateNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), prn)

	prs := openstack.TestPrivateSubnet(t)
	prs.NetworkID = prn.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), prs)

	f := openstack.TestFlavor(t)
	fetcher.Flavor().Create(openstack.TestRequestContext(t, fetcher, clusterID), f)

	sg := openstack.TestSecurityGroup(t)
	fetcher.SecurityGroup().Create(openstack.TestRequestContext(t, fetcher, clusterID), sg)

	kp := openstack.TestKeyPair(t)
	fetcher.KeyPair().Create(openstack.TestRequestContext(t, fetcher, clusterID), kp)

	pun := openstack.TestPublicNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pun)

	pus := openstack.TestPublicSubnet(t)
	pus.NetworkID = pun.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), pus)

	r := openstack.TestRouter(t)
	r.ExternalGatewayInfo.NetworkID = pun.ID
	fetcher.Router().Create(openstack.TestRequestContext(t, fetcher, clusterID), r)
	fetcher.Router().AddSubnet(openstack.TestRequestContext(t, fetcher, clusterID), r.ID, prs.ID)

	ip := &models.FloatingIp{}
	ip.NetworkID = pun.ID
	fetcher.FloatingIp().Create(openstack.TestRequestContext(t, fetcher, clusterID), ip)

	se := openstack.TestServer(t)
	se.FlavorID = f.ID
	se.SecurityGroupID = sg.ID
	se.KeyID = kp.ID
	se.PrivateNetworkID = prn.ID

	fetcher.Server().Create(openstack.TestRequestContext(t, fetcher, clusterID), se)

	time.Sleep(5000000000)

	ports, _ := fetcher.Port().FetchByDeviceID(openstack.TestRequestContext(t, fetcher, clusterID), se.ID)

	err := fetcher.FloatingIp().AddToPort(openstack.TestRequestContext(t, fetcher, clusterID), ip.ID, ports[0].ID)
	assert.NoError(t, err)

	fetcher.Server().Delete(openstack.TestRequestContext(t, fetcher, clusterID), se.ID)
	fetcher.Flavor().Delete(openstack.TestRequestContext(t, fetcher, clusterID), f.ID)
	fetcher.KeyPair().Delete(openstack.TestRequestContext(t, fetcher, clusterID), kp.ID)
	fetcher.SecurityGroup().Delete(openstack.TestRequestContext(t, fetcher, clusterID), sg.ID)
	fetcher.Router().RemoveExternalGateway(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Router().RemoveSubnet(openstack.TestRequestContext(t, fetcher, clusterID), r.ID, prs.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), prn.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), pun.ID)

}
