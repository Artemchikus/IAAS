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
		Name:        "test",
		DomainID:    "default",
		Description: "Test project",
	}
}

func TestImage(t *testing.T) *models.Image {
	return &models.Image{
		Name:            "Cirros",
		DiskFormat:      "qcow2",
		ContainerFormat: "bare",
		Visibility:      "public",
	}
}

func TestData(t *testing.T) []byte {
	file := "../../../test/cirros-0.6.1-x86_64-disk.img"

	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func TestFlavor(t *testing.T) *models.Flavor {
	return &models.Flavor{
		Disk:        2,
		RAM:         300,
		VCPUs:       1,
		Name:        "m1.small",
		Ephemeral:   0,
		IsPublic:    true,
		RXTXFactor:  1.0,
		Swap:        "0",
		Description: "test flavor",
	}
}

func TestPublicNetwork(t *testing.T) *models.Network {
	return &models.Network{
		Name:            "test-public",
		NetworkType:     "flat",
		External:        true,
		PhysicalNetwork: "external",
		MTU:             1442,
		Description:     "test network",
	}
}

func TestPrivateNetwork(t *testing.T) *models.Network {
	return &models.Network{
		Name:        "test-private",
		NetworkType: "geneve",
		External:    false,
		MTU:         1442,
		Description: "test network",
	}
}

func TestPublicSubnet(t *testing.T) *models.Subnet {
	allocationPool := &models.AllocationPool{
		Start: "192.168.122.100",
		End:   "192.168.122.110",
	}

	allocationPools := []*models.AllocationPool{allocationPool}

	return &models.Subnet{
		CIDR:            "192.168.122.0/24",
		Name:            "public-subnet",
		EnableDHCP:      false,
		AllocationPools: allocationPools,
		IpVersion:       4,
		GatewayIp:       "192.168.122.1",
		Description:     "test public subnet",
	}
}

func TestPrivateSubnet(t *testing.T) *models.Subnet {
	allocationPool := &models.AllocationPool{
		Start: "192.168.100.2",
		End:   "192.168.100.254",
	}
	allocationPools := []*models.AllocationPool{allocationPool}

	return &models.Subnet{
		CIDR:            "192.168.100.0/24",
		Name:            "private-subnet",
		EnableDHCP:      true,
		AllocationPools: allocationPools,
		IpVersion:       4,
		GatewayIp:       "192.168.100.1",
		Description:     "test private subnet",
	}
}

func TestRole(t *testing.T) *models.Role {
	return &models.Role{
		Name:        "test",
		Description: "Test role",
	}
}

func TestRouter(t *testing.T) *models.Router {
	info := &models.ExternalGatewayInfo{}

	return &models.Router{
		Description:         "Test router",
		Name:                "test",
		ExternalGatewayInfo: info,
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
		Bootable:    false,
		TypeID:      "9f315d13-84ad-451e-ad6f-26e673ff9144",
	}
}

func TestRequestContext(t *testing.T, fetcher business.Fetcher, clusterId int) context.Context {
	ctx := context.WithValue(context.Background(), models.CtxKeyClusterID, "99999999-9999-9999-9999-999999999999")
	ctx = context.WithValue(ctx, models.CtxKeyClusterID, clusterId)

	admin := models.TestClusters(t)[clusterId-1].Admin
	token, _ := fetcher.Token().Get(ctx, admin)
	return context.WithValue(ctx, models.CtxKeyToken, token)
}

func TestClusterUser(t *testing.T) *models.ClusterUser {
	return &models.ClusterUser{
		Email:       "test@example.com",
		Name:        "test",
		Password:    "password",
		DomainID:    "default",
		Description: "test user",
	}
}

func TestServer(t *testing.T) *models.Server {
	return &models.Server{
		Name:    "test",
		ImageID: "a043cd1d-8a15-48ea-a531-43aa68c41a20",
	}
}
