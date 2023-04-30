package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRouterFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	pubNet := openstack.TestPublicNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubNet)

	pubSub := openstack.TestPublicSubnet(t)
	pubSub.NetworkID = pubNet.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubSub)

	r := openstack.TestRouter(t)
	r.ExternalGatewayInfo.NetworkID = pubNet.ID

	err := fetcher.Router().Create(openstack.TestRequestContext(t, fetcher, clusterID), r)
	assert.NoError(t, err)
	assert.NotEqual(t, r.ID, "")

	time.Sleep(1000)

	fetcher.Router().RemoveExternalGateway(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Router().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ExternalGatewayInfo.NetworkID)
}

func TestRouterFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	pubNet := openstack.TestPublicNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubNet)

	pubSub := openstack.TestPublicSubnet(t)
	pubSub.NetworkID = pubNet.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubSub)

	r := openstack.TestRouter(t)
	r.ExternalGatewayInfo.NetworkID = pubNet.ID

	fetcher.Router().Create(openstack.TestRequestContext(t, fetcher, clusterID), r)

	time.Sleep(1000)

	fetcher.Router().RemoveExternalGateway(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)

	err := fetcher.Router().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, r.ID, "")
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ExternalGatewayInfo.NetworkID)
}

func TestRouterFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	pubNet := openstack.TestPublicNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubNet)

	pubSub := openstack.TestPublicSubnet(t)
	pubSub.NetworkID = pubNet.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubSub)

	r := openstack.TestRouter(t)
	r.ExternalGatewayInfo.NetworkID = pubNet.ID

	fetcher.Router().Create(openstack.TestRequestContext(t, fetcher, clusterID), r)

	time.Sleep(1000)

	r2, err := fetcher.Router().FetchByID(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, r2.ID, "")

	fetcher.Router().RemoveExternalGateway(openstack.TestRequestContext(t, fetcher, clusterID), r2.ID)
	fetcher.Router().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r2.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r2.ExternalGatewayInfo.NetworkID)
}

func TestRouterFetcher_FetchAll(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	pubNet := openstack.TestPublicNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubNet)

	pubSub := openstack.TestPublicSubnet(t)
	pubSub.NetworkID = pubNet.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubSub)

	r := openstack.TestRouter(t)
	r.ExternalGatewayInfo.NetworkID = pubNet.ID

	fetcher.Router().Create(openstack.TestRequestContext(t, fetcher, clusterID), r)

	time.Sleep(1000)

	rs, err := fetcher.Router().FetchAll(openstack.TestRequestContext(t, fetcher, clusterID))
	assert.NoError(t, err)
	assert.NotEmpty(t, rs)

	fetcher.Router().RemoveExternalGateway(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Router().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ExternalGatewayInfo.NetworkID)
}

func TestRouterFetcher_AddSubnet(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	pubNet := openstack.TestPublicNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubNet)

	pubSub := openstack.TestPublicSubnet(t)
	pubSub.NetworkID = pubNet.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubSub)

	r := openstack.TestRouter(t)
	r.ExternalGatewayInfo.NetworkID = pubNet.ID
	fetcher.Router().Create(openstack.TestRequestContext(t, fetcher, clusterID), r)

	privNet := openstack.TestPrivateNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), privNet)

	privSub := openstack.TestPrivateSubnet(t)
	privSub.NetworkID = privNet.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), privSub)

	time.Sleep(1000)

	err := fetcher.Router().AddSubnet(openstack.TestRequestContext(t, fetcher, clusterID), r.ID, privSub.ID)
	assert.NoError(t, err)

	fetcher.Router().RemoveExternalGateway(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Router().RemoveSubnet(openstack.TestRequestContext(t, fetcher, clusterID), r.ID, privSub.ID)
	fetcher.Router().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ExternalGatewayInfo.NetworkID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), privNet.ID)
}

func TestRouterFetcher_RemoveSubnet(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	pubNet := openstack.TestPublicNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubNet)

	pubSub := openstack.TestPublicSubnet(t)
	pubSub.NetworkID = pubNet.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubSub)

	r := openstack.TestRouter(t)
	r.ExternalGatewayInfo.NetworkID = pubNet.ID
	fetcher.Router().Create(openstack.TestRequestContext(t, fetcher, clusterID), r)

	privNet := openstack.TestPrivateNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), privNet)

	privSub := openstack.TestPrivateSubnet(t)
	privSub.NetworkID = privNet.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), privSub)

	fetcher.Router().AddSubnet(openstack.TestRequestContext(t, fetcher, clusterID), r.ID, privSub.ID)

	time.Sleep(1000)

	err := fetcher.Router().RemoveSubnet(openstack.TestRequestContext(t, fetcher, clusterID), r.ID, privSub.ID)
	assert.NoError(t, err)

	fetcher.Router().RemoveExternalGateway(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Router().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ExternalGatewayInfo.NetworkID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), privNet.ID)
}

func TestRouterFetcher_RemoveExternalGateway(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster", "clusterUser")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	pubNet := openstack.TestPublicNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubNet)

	pubSub := openstack.TestPublicSubnet(t)
	pubSub.NetworkID = pubNet.ID
	fetcher.Subnet().Create(openstack.TestRequestContext(t, fetcher, clusterID), pubSub)

	r := openstack.TestRouter(t)
	r.ExternalGatewayInfo.NetworkID = pubNet.ID

	fetcher.Router().Create(openstack.TestRequestContext(t, fetcher, clusterID), r)

	time.Sleep(1000)

	err := fetcher.Router().RemoveExternalGateway(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	assert.NoError(t, err)

	fetcher.Router().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ExternalGatewayInfo.NetworkID)
}
