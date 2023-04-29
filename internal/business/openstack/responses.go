package openstack

import (
	"IAAS/internal/models"
	"time"
)

type GetTokenResponse struct {
	Token *models.Token `json:"token"`
}

type CreateProjectResponse struct {
	Project *models.Project `json:"project"`
}

type FetchProjectResponse struct {
	Project *models.Project `json:"project"`
}

type CreateUserResponse struct {
	User *models.ClusterUser `json:"user"`
}

type FetchUserResponse struct {
	User *models.ClusterUser `json:"user"`
}

type FetchFlavorResponse struct {
	Flavor *models.Flavor `json:"flavor"`
}

type CreateFlavorResponse struct {
	Flavor *models.Flavor `json:"flavor"`
}

type FetchFloatingIpResponse struct {
	FloatingIp *models.FloatingIp `json:"floatingip"`
}

type CreateFloatingIpResponse struct {
	FloatingIp *models.FloatingIp `json:"floatingip"`
}

type FetchNetworkResponse struct {
	Network *models.Network `json:"network"`
}

type CreateNetworkResponse struct {
	Network *models.Network `json:"network"`
}

type CreateSubnetResponse struct {
	Subnet *models.Subnet `json:"subnet"`
}

type FetchSubnetResponse struct {
	Subnet *models.Subnet `json:"subnet"`
}

type CreateRoleResponse struct {
	Role *models.Role `json:"role"`
}

type FetchRoleResponse struct {
	Role *models.Role `json:"role"`
}

type CreateRouterResponse struct {
	Router *models.Router `json:"router"`
}

type FetchRouterResponse struct {
	Router *models.Router `json:"router"`
}

type CreateSecurityGroupResponse struct {
	SecurityGroup *models.SecurityGroup `json:"security_group"`
}

type FetchSecurityGroupResponse struct {
	SecurityGroup *models.SecurityGroup `json:"security_group"`
}

type FetchSecurityRuleResponse struct {
	SecurityRule *models.SecurityRule `json:"security_group_rule"`
}

type CreateSecurityRuleResponse struct {
	SecurityRule *models.SecurityRule `json:"security_group_rule"`
}

type FetchKeyPairResponse struct {
	KeyPair *map[string]interface{} `json:"keypair"`
}

type CreateKeyPairResponse struct {
	KeyPair *models.KeyPair `json:"keypair"`
}

type CreateVolumeResponse struct {
	Volume *map[string]interface{} `json:"volume"`
}

type FetchVolumeResponse struct {
	Volume *map[string]interface{} `json:"volume"`
}

type FetchPortsResponse struct {
	Ports []*models.Port `json:"ports"`
}

type FetchPortResponse struct {
	Port *models.Port `json:"port"`
}

type CreateServerResponse struct {
	Server *models.Server `json:"server"`
}

type FetchServerResponse struct {
	Server *FetchedServer `json:"server"`
}

type FetchedServer struct {
	HypervisorHostname string                `json:"OS-EXT-SRV-ATTR:hypervisor_hostname"`
	InstanceName       string                `json:"OS-EXT-SRV-ATTR:instance_name"`
	VMState            string                `json:"OS-EXT-STS:vm_state"`
	LaunchedAt         string                `json:"OS-SRV-USG:launched_at"`
	Addresses          map[string][]*Address `json:"addresses"`
	CreatedAt          time.Time             `json:"created"`
	Flavor             *IDResponse           `json:"flavor"`
	ID                 string                `json:"id"`
	Image              *IDResponse           `json:"image"`
	Key                string                `json:"key_name"`
	Name               string                `json:"name"`
	Volumes            []*IDResponse         `json:"os-extended-volumes:volumes_attached"`
	SecurityGroups     []*SGNameResponse     `json:"security_groups"`
	Status             string                `json:"status"`
	TenantID           string                `json:"tenant_id"`
	UpdatedAt          time.Time             `json:"updated"`
	UserID             string                `json:"user_id"`
}

type Address struct {
	Type    string `json:"OS-EXT-IPS:type"`
	Address string `json:"addr"`
}

type SGNameResponse struct {
	Name string `json:"name"`
}

type IDResponse struct {
	ID string `json:"id"`
}
type ErrorResponse struct {
	Error *Error `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Title   string `json:"title"`
}
