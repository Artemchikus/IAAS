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
