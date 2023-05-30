package openstack

import "IAAS/internal/models"

type GetTokenRequest struct {
	Auth *GetTokenAuth `json:"auth"`
}
type RefreshTokenRequest struct {
	Auth *RefreshTokenAuth `json:"auth"`
}
type RefreshTokenAuth struct {
	Identity *TokenIdentity `json:"identity"`
}
type TokenIdentity struct {
	Methods []string `json:"methods"`
	Token   *Token   `json:"password"`
}
type GetTokenAuth struct {
	Identity *PasswordIdentity `json:"identity"`
	Scope    *Scope            `json:"scope"`
}
type PasswordIdentity struct {
	Methods  []string  `json:"methods"`
	Password *Password `json:"password"`
}
type Scope struct {
	Project *GetTokenProject `json:"project"`
}
type GetTokenProject struct {
	ID     string  `json:"id"`
	Domain *Domain `json:"domain"`
}
type Password struct {
	User *User `json:"user"`
}
type User struct {
	Name     string  `json:"name"`
	Domain   *Domain `json:"domain"`
	Password string  `json:"password"`
}
type Domain struct {
	ID string `json:"id"`
}
type Token struct {
	ID string `json:"id"`
}
type CreateProjectRequest struct {
	Project *CreateProject `json:"project"`
}
type CreateProject struct {
	DomainID    string `json:"domain_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type CreateUserRequest struct {
	User *CreateUser `json:"user"`
}
type CreateUser struct {
	Name        string `json:"name"`
	DomainID    string `json:"domain_id"`
	ProjectID   string `json:"default_project_id"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Description string `json:"description"`
}
type CretaeImageRequest struct {
	DiskFormat      string `json:"disk_format"`
	ContainerFormat string `json:"container_format"`
	Name            string `json:"name"`
	Visibility      string `json:"visibility"`
}

type CreateFlavorRequest struct {
	Flavor *Flavor `json:"flavor"`
}

type Flavor struct {
	VCPUs      int     `json:"vcpus"`
	Disk       int     `json:"disk"`
	Name       string  `json:"name"`
	RAM        int     `json:"ram"`
	Ephemeral  int     `json:"OS-FLV-EXT-DATA:ephemeral"`
	IsPublic   bool    `json:"os-flavor-access:is_public"`
	Swap       string  `json:"swap"`
	RXTXFactor float32 `json:"rxtx_factor"`
	// Description string  `json:"description"`
}

type CreateFloatingIpRequest struct {
	FloatingIp *FloatingIp `json:"floatingip"`
}

type FloatingIp struct {
	NetworkID   string `json:"floating_network_id"`
	Description string `json:"description"`
}

type AddIpToPortRequest struct {
	FloatingIp *AddIpToPort `json:"floatingip"`
}

type AddIpToPort struct {
	PortID string `json:"port_id"`
}

type CreateNetworkRequest struct {
	Network *Network `json:"network"`
}

type Network struct {
	Name            string `json:"name"`
	NetworkType     string `json:"provider:network_type"`
	External        bool   `json:"router:external"`
	PhysicalNetwork string `json:"provider:physical_network"`
	MTU             int    `json:"mtu"`
	Description     string `json:"description"`
	ProjectID       string `json:"project_id"`
}

type CreateSubnetRequest struct {
	Subnet *Subnet `json:"subnet"`
}

type Subnet struct {
	CIDR            string                   `json:"cidr"`
	Name            string                   `json:"name"`
	EnableDHCP      bool                     `json:"enable_dhcp"`
	NetworkID       string                   `json:"network_id"`
	AllocationPools []*models.AllocationPool `json:"allocation_pools"`
	IpVersion       int                      `json:"ip_version"`
	GatewayIp       string                   `json:"gateway_ip"`
	Description     string                   `json:"description"`
}

type CreateRoleRequest struct {
	Role *Role `json:"role"`
}

type Role struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateRouterRequest struct {
	Router *Router `json:"router"`
}

type Router struct {
	Name                string                      `json:"name"`
	Description         string                      `json:"description"`
	ExternalGatewayInfo *models.ExternalGatewayInfo `json:"external_gateway_info"`
}

type RemoveRouterSubnetRequest struct {
	SubnetId string `json:"subnet_id"`
}

type RemoveRouterExternalGatewayRequest struct {
	Router *NullRouter `json:"router"`
}

type NullRouter struct {
	ExternalGatewayInfo *NullExternalInfo `json:"external_gateway_info"`
}

type NullExternalInfo struct {
}

type CreateSecurityGroupRequest struct {
	SecurityGroup *SecurityGroup `json:"security_group"`
}

type SecurityGroup struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateSecurityRuleRequest struct {
	SecurityRule *SecurityRule `json:"security_group_rule"`
}

type SecurityRule struct {
	Protocol        string `json:"protocol"`
	PortRangeMax    int    `json:"port_range_max"`
	RemoteIpPrefix  string `json:"remote_ip_prefix"`
	Ethertype       string `json:"ethertype"`
	SecurityGroupID string `json:"security_group_id"`
	Direction       string `json:"direction"`
	PortRangeMin    int    `json:"port_range_min"`
	Description     string `json:"description"`
}

type CreateKeyPairRequest struct {
	KeyPair *KeyPair `json:"keypair"`
}

type KeyPair struct {
	PublicKey string `json:"public_key"`
	Name      string `json:"name"`
	// Type      string `json:"type"`
}

type CreateVolumeRequest struct {
	Volume *Volume `json:"volume"`
}

type Volume struct {
	Size        int    `json:"size"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Bootable    bool   `json:"bootable"`
}

type CreateServerRequest struct {
	Server *Server `json:"server"`
}

type StartServerRequest struct {
	Start *string `json:"os-start"`
}

type StopServerRequest struct {
	Stop *string `json:"os-stop"`
}

type AddVolumeToServerRequest struct {
	VolumeAttachment *VolumeAttachment `json:"volumeAttachment"`
}

type VolumeAttachment struct {
	VolumeId string  `json:"volumeId"`
	Device   *string `json:"device"`
}

type Server struct {
	FlavorID       string                   `json:"flavorRef"`
	ImageID        string                   `json:"imageRef"`
	KeyID          string                   `json:"key_name"`
	Name           string                   `json:"name"`
	SecurityGroups []*ServerSecurityGroupID `json:"security_groups"`
	Networks       []*ServerNetworkID       `json:"networks"`
}

type ServerSecurityGroupID struct {
	ID string `json:"name"`
}

type ServerNetworkID struct {
	ID string `json:"uuid"`
}
