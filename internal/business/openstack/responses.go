package openstack

import "IAAS/internal/models"

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
	User *map[string]interface{} `json:"user"`
}

type FetchUserResponse struct {
	User *map[string]interface{} `json:"user"`
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

type ErrorResponse struct {
	Error *Error `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Title   string `json:"title"`
}
