package openstack

import (
	"IAAS/internal/config"
	"IAAS/internal/models"
	"testing"
)

func TestConfig(t *testing.T) *config.ApiConfig {
	config := config.NewConfig()
	config.JwtKey = "secretkey"
	config.Admin = models.TestAdmin(t)
	config.Clusters = TestClusters(t)

	return config
}

func TestClusters(t *testing.T) []*models.Cluster {
	cluster := &models.Cluster{
		ID:       1,
		Location: "rus",
		URL:      "http://192.168.122.20:5000",
		Admin:    TestClusterAdmin(t),
	}

	return []*models.Cluster{cluster}
}

func TestCluster(t *testing.T) *models.Cluster {
	return &models.Cluster{
		Location: "test",
		URL:      "test",
		Admin:    TestClusterAdmin(t),
	}
}

func TestClusterAdmin(t *testing.T) *models.Account {
	return &models.Account{
		Email:    "adm@example.com",
		Name:     "admin",
		Password: "openstack",
	}
}

func TestProject(t *testing.T) *models.Project {
	return &models.Project{
		Name:        "demo",
		Enabled:     true,
		DomainID:    "default",
		Description: "Demo project",
		Options:     &models.Options{},
		Tags:        make([]string, 0),
	}
}
