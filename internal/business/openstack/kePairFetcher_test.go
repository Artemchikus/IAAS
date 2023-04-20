package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKeyPairFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	kp := openstack.TestKeyPair(t)

	err := fetcher.KeyPair().Create(models.TestRequestContext(t), clusterID, kp)
	assert.NoError(t, err)
	assert.NotEqual(t, kp.Fingerprint, "")

	time.Sleep(1000)

	fetcher.KeyPair().Delete(models.TestRequestContext(t), clusterID, kp.Name)
}

func TestKeyPairFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	kp := openstack.TestKeyPair(t)

	fetcher.KeyPair().Create(models.TestRequestContext(t), clusterID, kp)

	time.Sleep(1000)

	err := fetcher.KeyPair().Delete(models.TestRequestContext(t), clusterID, kp.Name)
	assert.NoError(t, err)
	assert.NotEqual(t, kp.Name, "")
}

func TestKeyPairFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	kp1 := openstack.TestKeyPair(t)

	fetcher.KeyPair().Create(models.TestRequestContext(t), clusterID, kp1)

	time.Sleep(1000)

	kp2, err := fetcher.KeyPair().FetchByID(models.TestRequestContext(t), clusterID, kp1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, kp2.Fingerprint, "")

	fetcher.KeyPair().Delete(models.TestRequestContext(t), clusterID, kp2.Name)
}
