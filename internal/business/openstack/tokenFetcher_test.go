package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const databaseURL = "host=localhost port=5433 user=postgres password=iaas dbname=iaas-test sslmode=disable"

func TestTokenFetcher_Get(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.NewStore(models.TestInitContext(t), db, config)

	fetcher := openstack.NewFetcher(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	account := config.Clusters[0].Admin

	token, err := fetcher.Token().Get(openstack.TestRequestContext(t, fetcher, clusterID), account)
	assert.NoError(t, err)
	assert.NotNil(t, token.Value)
	assert.NotNil(t, token.ExpiresAt)
	var zeroTime time.Time
	assert.NotEqual(t, token.ExpiresAt, zeroTime)
}
