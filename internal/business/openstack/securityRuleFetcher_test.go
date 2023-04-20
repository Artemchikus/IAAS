package openstack_test

import (
	"IAAS/internal/business/openstack"
	"IAAS/internal/models"
	"IAAS/internal/store/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSecurityRuleFetcher_Create(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	sg := openstack.TestSecurityGroup(t)

	fetcher.SecurityGroup().Create(models.TestRequestContext(t), clusterID, sg)

	sr := openstack.TestSecurityRule(t)
	sr.SecurityGroupID = sg.ID

	err := fetcher.SecurityRule().Create(models.TestRequestContext(t), clusterID, sr)
	assert.NoError(t, err)
	assert.NotEqual(t, sr.ID, "")

	time.Sleep(1000)

	fetcher.SecurityGroup().Delete(models.TestRequestContext(t), clusterID, sg.ID)
}

func TestSecurityRuleFetcher_Delete(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	sg := openstack.TestSecurityGroup(t)

	fetcher.SecurityGroup().Create(models.TestRequestContext(t), clusterID, sg)

	sr := openstack.TestSecurityRule(t)
	sr.SecurityGroupID = sg.ID

	fetcher.SecurityRule().Create(models.TestRequestContext(t), clusterID, sr)

	time.Sleep(1000)

	err := fetcher.SecurityRule().Delete(models.TestRequestContext(t), clusterID, sr.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, sr.ID, "")

	fetcher.SecurityGroup().Delete(models.TestRequestContext(t), clusterID, sg.ID)
}

func TestSecurityRuleFetcher_FetchByID(t *testing.T) {
	db, teardown := postgres.TestDB(t, databaseURL)
	defer teardown("account", "secret", "cluster")

	config := openstack.TestConfig(t)

	s := postgres.New(models.TestInitContext(t), db, config)

	fetcher := openstack.New(models.TestInitContext(t), config, s)

	clusterID := config.Clusters[0].ID

	sg := openstack.TestSecurityGroup(t)

	fetcher.SecurityGroup().Create(models.TestRequestContext(t), clusterID, sg)

	sr1 := openstack.TestSecurityRule(t)
	sr1.SecurityGroupID = sg.ID

	fetcher.SecurityRule().Create(models.TestRequestContext(t), clusterID, sr1)

	time.Sleep(1000)

	sr2, err := fetcher.SecurityRule().FetchByID(models.TestRequestContext(t), clusterID, sr1.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, sr2.ID, "")

	fetcher.SecurityGroup().Delete(models.TestRequestContext(t), clusterID, sg.ID)
}
