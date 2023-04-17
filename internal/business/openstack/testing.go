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
	config.Clusters = models.TestClusters(t)

	return config
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
