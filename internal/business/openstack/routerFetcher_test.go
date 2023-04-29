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

	pn := openstack.TestPublicNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pn)

	r := openstack.TestRouter(t)
	r.ExternalGatewayInfo.NetworkID = pn.ID

	err := fetcher.Router().Create(openstack.TestRequestContext(t, fetcher, clusterID), r)
	assert.NoError(t, err)
	assert.NotEqual(t, r.ID, "")

	time.Sleep(1000)

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

	pn := openstack.TestPublicNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pn)

	r := openstack.TestRouter(t)
	r.ExternalGatewayInfo.NetworkID = pn.ID

	fetcher.Router().Create(openstack.TestRequestContext(t, fetcher, clusterID), r)

	time.Sleep(1000)

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

	pn := openstack.TestPublicNetwork(t)
	fetcher.Network().Create(openstack.TestRequestContext(t, fetcher, clusterID), pn)

	r1 := openstack.TestRouter(t)
	r1.ExternalGatewayInfo.NetworkID = pn.ID

	fetcher.Router().Create(openstack.TestRequestContext(t, fetcher, clusterID), r1)

	time.Sleep(1000)

	r2, err := fetcher.Router().FetchByID(openstack.TestRequestContext(t, fetcher, clusterID), r1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, r2.ID, "")

	fetcher.Router().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r2.ID)
	fetcher.Network().Delete(openstack.TestRequestContext(t, fetcher, clusterID), r2.ExternalGatewayInfo.NetworkID)
}
