package openstack

import (
	"IAAS/internal/business"
	"IAAS/internal/config"
	"IAAS/internal/models"
	"context"
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

func TestNetwork(t *testing.T) *models.Network {
	return &models.Network{
		Name:            "public",
		NetworkType:     "flat",
		AdminStateUp:    true,
		External:        true,
		PhysicalNetwork: "external",
	}
}

func TestSubnet(t *testing.T) *models.Subnet {
	allocationPool := &models.AllocationPool{
		Start: "192.168.122.200",
		End:   "192.168.122.254",
	}

	allocationPools := []*models.AllocationPool{allocationPool}

	return &models.Subnet{
		CIDR:            "192.168.122.0/24",
		Name:            "public-subnet",
		EnableDHCP:      false,
		AllocationPools: allocationPools,
		IpVersion:       4,
		GatewayIp:       "192.168.122.1",
	}
}

func TestRole(t *testing.T) *models.Role {
	return &models.Role{
		Name:        "test",
		Description: "Test role",
	}
}

func TestRouter(t *testing.T) *models.Router {
	return &models.Router{
		Description: "Test router",
		Name:        "test",
	}
}

func TestSecurityGroup(t *testing.T) *models.SecurityGroup {
	return &models.SecurityGroup{
		Description: "Test security group",
		Name:        "test",
	}
}

func TestSecurityRule(t *testing.T) *models.SecurityRule {
	return &models.SecurityRule{
		Protocol:       "tcp",
		Description:    "SSH security rule",
		PortRangeMax:   22,
		PortRangeMin:   22,
		Direction:      "ingress",
		RemoteIpPrefix: "0.0.0.0/0",
		Ethertype:      "IPv4",
	}
}

func TestKeyPair(t *testing.T) *models.KeyPair {
	return &models.KeyPair{
		Name:      "testKeyPair",
		PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDZjh2NAkpTS2QVB57w5p9vAU8WbmgjdrYwAklFEErMui+ZaBSWZUFHrgDS7NyTbLFF7VLXZHeKFrkF1YKYognlxvkCWyAqOPXFlR4Ajm8XF54RGBACpL+iWr8IodhguHC01Ygj23L9k5mNQ3ZLWHYh8zXACh5kCjc7Tu1yLPWnd4AvntmPv0ZfbE4xRpTQZI579EFZxmBVZdwJgql6lH/bFgktlYMJnXMyi+oAgaVSyLUwWCT9UHnKv6A/JvcrSIdI8enZwtKJPyNP4+LyDVkBhDPMJ7O4iR5iMK49bcCzT0W2ycuI1q5ZYMD4RN1CGcKaHO36bZJ4etfEchC+F98XUe2CNfXn61kRrL0N7qahvKa62ln9tHZuLGTwxTvxoHZ9ZNVVqsbsrOkgsJZfbRh7paD5FBAotaUv1TFtVyyeRLMRvlNVPeIclbRK5hur7ptryd25iS9IGbqnJ6VvSMs6hGCcfopEP28J5eN0gmx8pCsl2v0GkAqEILjWpl66cVc= root@controller.test.local\n",
		Type:      "ssh",
	}
}

func TestVolume(t *testing.T) *models.Volume {
	return &models.Volume{
		Name:        "test",
		Description: "test volume",
		Size:        1,
	}
}

func TestRequestContext(t *testing.T, fetcher business.Fetcher, clusterId int) context.Context {
	ctx := context.WithValue(context.Background(), models.CtxKeyClusterID, "99999999-9999-9999-9999-999999999999")
	ctx = context.WithValue(ctx, models.CtxKeyClusterID, clusterId)

	admin := models.TestClusters(t)[clusterId].Admin
	token, _ := fetcher.Token().Get(ctx, admin)
	return context.WithValue(ctx, models.CtxKeyToken, token)
}
