package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server struct {
		Port       string
		ApiVersion string
	}

	CrudApi struct {
		Port string
	}

	DefaultDecision bool
}

func NewConfig() *Config {
	cfg := &Config{}

	cfg.Server.Port = os.Getenv("DECISION_SERVER_PORT")
	cfg.Server.ApiVersion = os.Getenv("API_VERSION")

	cfg.CrudApi.Port = os.Getenv("CRUD_SERVER_PORT")

	defaultDecision, err := strconv.ParseBool(os.Getenv("DEFAULT_DECISION"))
	if err != nil {
		defaultDecision = false
	}
	cfg.DefaultDecision = defaultDecision

	return cfg
}
