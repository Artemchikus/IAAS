package config

import "IAAS/internal/models"

type ApiConfig struct {
	BindAddr    string            `toml:"bind_addr"`
	LogLevel    string            `toml:"log_level"`
	DatabaseURL string            `toml:"database_url"`
	JwtKey      string            `toml:"jwt_key"`
	Clusters    []*models.Cluster `toml:"clusters"`
	Admin       *models.Account   `toml:"api_admin"`
}

func NewConfig() *ApiConfig {
	return &ApiConfig{
		BindAddr: ":8080",
		LogLevel: "debug",
		Admin:    new(models.Account),
	}
}
