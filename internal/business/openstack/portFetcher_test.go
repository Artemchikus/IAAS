package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPortFetcher_FetchByID(t *testing.T) {
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

	ps, _ := fetcher.Port().FetchByNetworkID(openstack.TestRequestContext(t, fetcher, clusterID), pubNet.ID)

	p, err := fetcher.Port().FetchByID(openstack.TestRequestContext(t, fetcher, clusterID), ps[0].ID)
	assert.NoError(t, err)
	assert.NotEqual(t, p.ID, "")

	fetcher.Router().RemoveExternalGateway(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Router().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ExternalGatewayInfo.NetworkID)
}

func TestPortFetcher_FetchByRouterID(t *testing.T) {
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

	ps, err := fetcher.Port().FetchByRouterID(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, ps)

	fetcher.Router().RemoveExternalGateway(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Router().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ExternalGatewayInfo.NetworkID)
}

func TestPortFetcher_FetchByNetworkID(t *testing.T) {
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

	ps, err := fetcher.Port().FetchByNetworkID(openstack.TestRequestContext(t, fetcher, clusterID), pubNet.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, ps)

	fetcher.Router().RemoveExternalGateway(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Router().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r.ExternalGatewayInfo.NetworkID)
}
