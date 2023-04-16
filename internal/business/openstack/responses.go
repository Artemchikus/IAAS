package openstack

import "IAAS/internal/models"

type GetTokenResponse struct {
	Token *models.Token `json:"token"`
}

type CreateProjectResponse struct {
	Project *models.Project `json:"project"`
}

type findProjectResponse struct {
	Project *models.Project `json:"project"`
}
