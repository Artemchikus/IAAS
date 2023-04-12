package config

type ApiConfig struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	JwtKey      string `toml:"jwt_key"`
}

func NewConfig() *ApiConfig {
	return &ApiConfig{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}
