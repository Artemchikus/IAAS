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

	sub := openstack.TestSubnet(t)
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

	time.Sleep(1000)

	fetcher.Server().Delete(openstack.TestRequestContext(t, fetcher, clusterID), se.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), pn.ID)
	fetcher.Flavor().Delete(openstack.TestRequestContext(t, fetcher, clusterID), f.ID)
	fetcher.KeyPair().Delete(openstack.TestRequestContext(t, fetcher, clusterID), kp.ID)
	fetcher.SecurityGroup().Delete(openstack.TestRequestContext(t, fetcher, clusterID), sg.ID)
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

	sub := openstack.TestSubnet(t)
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

	time.Sleep(1000)

	err := fetcher.Server().Delete(openstack.TestRequestContext(t, fetcher, clusterID), se.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, se.ID, "")

	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), pn.ID)
	fetcher.Flavor().Delete(openstack.TestRequestContext(t, fetcher, clusterID), f.ID)
	fetcher.KeyPair().Delete(openstack.TestRequestContext(t, fetcher, clusterID), kp.ID)
	fetcher.SecurityGroup().Delete(openstack.TestRequestContext(t, fetcher, clusterID), sg.ID)
}

func TestServerFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	// pn := openstack.TestPrivateNetwork(t)
	// fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pn)

	// sub := openstack.TestSubnet(t)
	// sub.NetworkID = pn.ID
	// fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), sub)

	// f := openstack.TestFlavor(t)
	// fetcher.Flavor().Create(openstack.TestRequestContext(t, fetcher, clusterID), f)

	// sg := openstack.TestSecurityGroup(t)
	// fetcher.SecurityGroup().Create(openstack.TestRequestContext(t, fetcher, clusterID), sg)

	// kp := openstack.TestKeyPair(t)
	// fetcher.KeyPair().Create(openstack.TestRequestContext(t, fetcher, clusterID), kp)

	// se1 := openstack.TestServer(t)
	// se1.FlavorID = f.ID

	// sgID := &models.ServerSecurityGroupID{
	// 	ID: sg.ID,
	// }
	// se1.SecurityGroupsID = append(se1.SecurityGroupsID, sgID)

	// se1.KeyID = kp.ID

	// nId := &models.ServerNetworkID{
	// 	ID: pn.ID,
	// }
	// se1.NetworksID = append(se1.NetworksID, nId)

	// fetcher.Server().Create(openstack.TestRequestContext(t, fetcher, clusterID), se1)

	// time.Sleep(1000)

	se2, err := fetcher.Server().FetchByID(openstack.TestRequestContext(t, fetcher, clusterID), "d33d6eaa-51f4-48ff-84fd-6da7a1e7c466")
	assert.NoError(t, err)
	assert.NotEqual(t, se2.ID, "")

	// fetcher.Server().Delete(openstack.TestRequestContext(t, fetcher, clusterID), se2.ID)
	// fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), pn.ID)
	// fetcher.Flavor().Delete(openstack.TestRequestContext(t, fetcher, clusterID), f.ID)
	// fetcher.KeyPair().Delete(openstack.TestRequestContext(t, fetcher, clusterID), kp.ID)
	// fetcher.SecurityGroup().Delete(openstack.TestRequestContext(t, fetcher, clusterID), sg.ID)
}
