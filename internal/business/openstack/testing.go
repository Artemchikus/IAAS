package openstack

import (
	"IAAS/internal/config"
	"IAAS/internal/models"
	"log"
	"os"
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

func TestImage(t *testing.T) *models.Image {
	file := "../../../test/cirros-0.6.1-x86_64-disk.img"

	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(len(data))

	return &models.Image{
		FileData:        data,
		Name:            "Cirros",
		DiskFormat:      "qcow2",
		ContainerFormat: "bare",
		Visibility:      "public",
	}
}

func TestFlavor(t *testing.T) *models.Flavor {
	return &models.Flavor{
		Disk:       2,
		RAM:        300,
		VCPUs:      1,
		Name:       "m1.small",
		Ephemeral:  0,
		IsPublic:   true,
		RXTXFactor: 1.0,
		Swap:       "0",
	}
}
