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