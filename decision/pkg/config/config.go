package config

import "os"

type Config struct {
	Server struct {
		Port       string
		ApiVersion string
	}

	CrudApi struct {
		Port string
	}
}

func NewConfig() *Config {
	cfg := &Config{}

	cfg.Server.Port = os.Getenv("DECISION_SERVER_PORT")
	cfg.Server.ApiVersion = os.Getenv("API_VERSION")

	cfg.CrudApi.Port = os.Getenv("CRUD_SERVER_PORT")

	return cfg
}
