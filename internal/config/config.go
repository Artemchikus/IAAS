package config

import "IAAS/internal/models"

type ApiConfig struct {
	BindAddr    string          `toml:"bind_addr"`
	LogLevel    string          `toml:"log_level"`
	DatabaseURL string          `toml:"database_url"`
	JwtKey      string          `toml:"jwt_key"`
	Admin       *models.Account `toml:"admin"`
}

func NewConfig() *ApiConfig {
	return &ApiConfig{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}
