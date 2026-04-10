package config

import "os"

type Config struct {
	Server struct {
		Port string
	}

	JWT struct {
		SecretKey string
	}

	Database struct {
		Host     string
		User     string
		Password string
		Name     string
		Port     string
	}
}

func NewConfig() *Config {
	cfg := &Config{}

	cfg.Server.Port = os.Getenv("SERVER_PORT")

	cfg.JWT.SecretKey = os.Getenv("JWT_SECRET_KEY")

	cfg.Database.Host = os.Getenv("DB_HOST")
	cfg.Database.Port = os.Getenv("DB_PORT")
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.Name = os.Getenv("DB_NAME")

	return cfg
}
