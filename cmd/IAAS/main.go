package main

import (
	"IAAS/internal/api"
	"IAAS/internal/config"
	"flag"
	"log"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config_path", "./config/api.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := config.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := api.Start(config); err != nil {
		log.Fatal(err)
	}
}
