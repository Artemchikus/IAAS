package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	pn := openstack.TestPrivateNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pn)

	sub := openstack.TestPrivateSubnet(t)
	sub.NetworkID = pn.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), sub)

	f := openstack.TestFlavor(t)
	fetcher.Flavor().Create(openstack.TestRequestContext(t, fetcher, clusterID), f)

	sg := openstack.TestSecurityGroup(t)
	fetcher.SecurityGroup().Create(openstack.TestRequestContext(t, fetcher, clusterID), sg)

	kp := openstack.TestKeyPair(t)
	fetcher.KeyPair().Create(openstack.TestRequestContext(t, fetcher, clusterID), kp)

	se := openstack.TestServer(t)
	se.FlavorID = f.ID
	se.SecurityGroupID = sg.ID
	se.KeyID = kp.ID
	se.PrivateNetworkID = pn.ID

	err := fetcher.Server().Create(openstack.TestRequestContext(t, fetcher, clusterID), se)
	assert.NoError(t, err)
	assert.NotEqual(t, se.ID, "")

	time.Sleep(2000)

	fetcher.Server().Delete(openstack.TestRequestContext(t, fetcher, clusterID), se.ID)
	fetcher.Flavor().Delete(openstack.TestRequestContext(t, fetcher, clusterID), f.ID)
	fetcher.KeyPair().Delete(openstack.TestRequestContext(t, fetcher, clusterID), kp.ID)
	fetcher.SecurityGroup().Delete(openstack.TestRequestContext(t, fetcher, clusterID), sg.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), pn.ID)
}

func TestServerFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	pn := openstack.TestPrivateNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pn)

	sub := openstack.TestPrivateSubnet(t)
	sub.NetworkID = pn.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), sub)

	f := openstack.TestFlavor(t)
	fetcher.Flavor().Create(openstack.TestRequestContext(t, fetcher, clusterID), f)

	sg := openstack.TestSecurityGroup(t)
	fetcher.SecurityGroup().Create(openstack.TestRequestContext(t, fetcher, clusterID), sg)

	kp := openstack.TestKeyPair(t)
	fetcher.KeyPair().Create(openstack.TestRequestContext(t, fetcher, clusterID), kp)

	se := openstack.TestServer(t)
	se.FlavorID = f.ID
	se.SecurityGroupID = sg.ID
	se.KeyID = kp.ID
	se.PrivateNetworkID = pn.ID

	fetcher.Server().Create(openstack.TestRequestContext(t, fetcher, clusterID), se)

	time.Sleep(2000)

	err := fetcher.Server().Delete(openstack.TestRequestContext(t, fetcher, clusterID), se.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, se.ID, "")

	fetcher.Flavor().Delete(openstack.TestRequestContext(t, fetcher, clusterID), f.ID)
	fetcher.KeyPair().Delete(openstack.TestRequestContext(t, fetcher, clusterID), kp.ID)
	fetcher.SecurityGroup().Delete(openstack.TestRequestContext(t, fetcher, clusterID), sg.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), pn.ID)
}

func TestServerFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	pn := openstack.TestPrivateNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pn)

	sub := openstack.TestPrivateSubnet(t)
	sub.NetworkID = pn.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), sub)

	f := openstack.TestFlavor(t)
	fetcher.Flavor().Create(openstack.TestRequestContext(t, fetcher, clusterID), f)

	sg := openstack.TestSecurityGroup(t)
	fetcher.SecurityGroup().Create(openstack.TestRequestContext(t, fetcher, clusterID), sg)

	kp := openstack.TestKeyPair(t)
	fetcher.KeyPair().Create(openstack.TestRequestContext(t, fetcher, clusterID), kp)

	se1 := openstack.TestServer(t)
	se1.FlavorID = f.ID
	se1.SecurityGroupID = sg.ID
	se1.KeyID = kp.ID
	se1.PrivateNetworkID = pn.ID

	fetcher.Server().Create(openstack.TestRequestContext(t, fetcher, clusterID), se1)

	time.Sleep(5000000000)

	se2, err := fetcher.Server().FetchByID(openstack.TestRequestContext(t, fetcher, clusterID), se1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, se2.ID, "")

	fetcher.Server().Delete(openstack.TestRequestContext(t, fetcher, clusterID), se2.ID)
	fetcher.Flavor().Delete(openstack.TestRequestContext(t, fetcher, clusterID), f.ID)
	fetcher.KeyPair().Delete(openstack.TestRequestContext(t, fetcher, clusterID), kp.ID)
	fetcher.SecurityGroup().Delete(openstack.TestRequestContext(t, fetcher, clusterID), sg.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), pn.ID)
}

func TestServerFetcher_FetchAll(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	pn := openstack.TestPrivateNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pn)

	sub := openstack.TestPrivateSubnet(t)
	sub.NetworkID = pn.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), sub)

	f := openstack.TestFlavor(t)
	fetcher.Flavor().Create(openstack.TestRequestContext(t, fetcher, clusterID), f)

	sg := openstack.TestSecurityGroup(t)
	fetcher.SecurityGroup().Create(openstack.TestRequestContext(t, fetcher, clusterID), sg)

	kp := openstack.TestKeyPair(t)
	fetcher.KeyPair().Create(openstack.TestRequestContext(t, fetcher, clusterID), kp)

	se := openstack.TestServer(t)
	se.FlavorID = f.ID
	se.SecurityGroupID = sg.ID
	se.KeyID = kp.ID
	se.PrivateNetworkID = pn.ID

	fetcher.Server().Create(openstack.TestRequestContext(t, fetcher, clusterID), se)

	time.Sleep(5000000000)

	ses, err := fetcher.Server().FetchAll(openstack.TestRequestContext(t, fetcher, clusterID))
	assert.NoError(t, err)
	assert.NotEmpty(t, ses)

	fetcher.Server().Delete(openstack.TestRequestContext(t, fetcher, clusterID), se.ID)
	fetcher.Flavor().Delete(openstack.TestRequestContext(t, fetcher, clusterID), f.ID)
	fetcher.KeyPair().Delete(openstack.TestRequestContext(t, fetcher, clusterID), kp.ID)
	fetcher.SecurityGroup().Delete(openstack.TestRequestContext(t, fetcher, clusterID), sg.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), pn.ID)
}

func TestServerFetcher_StartAndStop(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	pn := openstack.TestPrivateNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pn)

	sub := openstack.TestPrivateSubnet(t)
	sub.NetworkID = pn.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), sub)

	f := openstack.TestFlavor(t)
	fetcher.Flavor().Create(openstack.TestRequestContext(t, fetcher, clusterID), f)

	sg := openstack.TestSecurityGroup(t)
	fetcher.SecurityGroup().Create(openstack.TestRequestContext(t, fetcher, clusterID), sg)

	kp := openstack.TestKeyPair(t)
	fetcher.KeyPair().Create(openstack.TestRequestContext(t, fetcher, clusterID), kp)

	se := openstack.TestServer(t)
	se.FlavorID = f.ID
	se.SecurityGroupID = sg.ID
	se.KeyID = kp.ID
	se.PrivateNetworkID = pn.ID

	fetcher.Server().Create(openstack.TestRequestContext(t, fetcher, clusterID), se)

	time.Sleep(5000000000)

	err := fetcher.Server().Stop(openstack.TestRequestContext(t, fetcher, clusterID), se.ID)
	assert.NoError(t, err)

	time.Sleep(5000000000)

	err = fetcher.Server().Start(openstack.TestRequestContext(t, fetcher, clusterID), se.ID)
	assert.NoError(t, err)

	time.Sleep(5000000000)

	fetcher.Server().Delete(openstack.TestRequestContext(t, fetcher, clusterID), se.ID)
	fetcher.Flavor().Delete(openstack.TestRequestContext(t, fetcher, clusterID), f.ID)
	fetcher.KeyPair().Delete(openstack.TestRequestContext(t, fetcher, clusterID), kp.ID)
	fetcher.SecurityGroup().Delete(openstack.TestRequestContext(t, fetcher, clusterID), sg.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), pn.ID)
}

func TestServerFetcher_AttachVolume(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	pn := openstack.TestPrivateNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pn)

	sub := openstack.TestPrivateSubnet(t)
	sub.NetworkID = pn.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), sub)

	f := openstack.TestFlavor(t)
	fetcher.Flavor().Create(openstack.TestRequestContext(t, fetcher, clusterID), f)

	sg := openstack.TestSecurityGroup(t)
	fetcher.SecurityGroup().Create(openstack.TestRequestContext(t, fetcher, clusterID), sg)

	kp := openstack.TestKeyPair(t)
	fetcher.KeyPair().Create(openstack.TestRequestContext(t, fetcher, clusterID), kp)

	v := openstack.TestVolume(t)
	fetcher.Volume().Create(openstack.TestRequestContext(t, fetcher, clusterID), v)

	se := openstack.TestServer(t)
	se.FlavorID = f.ID
	se.SecurityGroupID = sg.ID
	se.KeyID = kp.ID
	se.PrivateNetworkID = pn.ID

	fetcher.Server().Create(openstack.TestRequestContext(t, fetcher, clusterID), se)

	time.Sleep(5000000000)

	err := fetcher.Server().AttachVolume(openstack.TestRequestContext(t, fetcher, clusterID), se.ID, v.ID)
	assert.NoError(t, err)

	fetcher.Server().Delete(openstack.TestRequestContext(t, fetcher, clusterID), se.ID)
	fetcher.Flavor().Delete(openstack.TestRequestContext(t, fetcher, clusterID), f.ID)
	fetcher.KeyPair().Delete(openstack.TestRequestContext(t, fetcher, clusterID), kp.ID)
	fetcher.SecurityGroup().Delete(openstack.TestRequestContext(t, fetcher, clusterID), sg.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), pn.ID)
	fetcher.Volume().Delete(openstack.TestRequestContext(t, fetcher, clusterID), v.ID)
}
