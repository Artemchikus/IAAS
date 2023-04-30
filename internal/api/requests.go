package api

import "IAAS/internal/models"

type CreateAccountRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateServerRequest struct {
	ImageID          string `json:"image_id"`
	KeyID            string `json:"key_id"`
	Name             string `json:"name"`
	SecurityGroupID  string `json:"security_groups"`
	PrivateNetworkID string `json:"private_network_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateClusterRequest struct {
	Location string                     `json:"location"`
	Url      string                     `json:"url"`
	Admin    *CreateClusterAdminRequest `json:"admin"`
}

type CreateFlavorRequest struct {
	VCPUs       int     `json:"vcpus"`
	Disk        int     `json:"disk"`
	Name        string  `json:"name"`
	RAM         int     `json:"ram"`
	Ephemeral   int     `json:"ephemeral"`
	IsPublic    bool    `json:"is_public"`
	Swap        int     `json:"swap"`
	RXTXFactor  float32 `json:"rxtx_factor"`
	Description string  `json:"description"`
}

type CreateFloatingIpRequest struct {
	NetworkID   string `json:"floating_network_id"`
	Description string `json:"description"`
}

type CreateImageRequest struct {
	DiskFormat      string `json:"disk_format"`
	ContainerFormat string `json:"container_format"`
	Name            string `json:"name"`
	Visibility      string `json:"visibility"`
}

type CreateKeyPairRequest struct {
	PublicKey string `json:"public_key"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}

type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type CreateNetworkRequest struct {
	Name            string `json:"name"`
	NetworkType     string `json:"network_type"`
	External        bool   `json:"is_external"`
	PhysicalNetwork string `json:"physical_network"`
	MTU             int    `json:"mtu"`
	Description     string `json:"description"`
}

type CreateProjectRequest struct {
	DomainID    string `json:"domain_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateRouterRequest struct {
	Name                string                      `json:"name"`
	Description         string                      `json:"description"`
	ExternalGatewayInfo *models.ExternalGatewayInfo `json:"external_gateway_info"`
}

type CreateSecurityGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateSecurityRuleRequest struct {
	Protocol        string `json:"protocol"`
	PortRangeMax    int    `json:"port_range_max"`
	RemoteIpPrefix  string `json:"remote_ip_prefix"`
	Ethertype       string `json:"ethertype"`
	SecurityGroupID string `json:"security_group_id"`
	Direction       string `json:"direction"`
	PortRangeMin    int    `json:"port_range_min"`
	Description     string `json:"description"`
}

type CreateSubnetRequest struct {
	CIDR            string                   `json:"cidr"`
	Name            string                   `json:"name"`
	EnableDHCP      bool                     `json:"enable_dhcp"`
	NetworkID       string                   `json:"network_id"`
	AllocationPools []*models.AllocationPool `json:"allocation_pools"`
	IpVersion       int                      `json:"ip_version"`
	GatewayIp       string                   `json:"gateway_ip"`
	Description     string                   `json:"description"`
}

type CreateUserRequest struct {
	Name        string `json:"name"`
	ProjectID   string `json:"project_id"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	DomainID    string `json:"domain_id"`
	Description string `json:"description"`
}

type CreateClusterAdminRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ProjectID string `json:"project_id"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	DomainID  string `json:"domain_id"`
}

type CreateVolumeRequest struct {
	Size        int    `json:"size"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TypeID      string `json:"volume_type"`
	Bootable    bool   `json:"bootable"`
}

type AttachVolumeRequest struct {
	VolumeID string `json:"volume_id"`
}

type RemoveOrAddSubnetRequest struct {
	SubnetID string `json:"subnet_id"`
}

type AddFloatingIPRequest struct {
	PortID string `json:"port_id"`
}
