package config

import "os"

type Config struct {
	Server struct {
		Address string
	}
}

func LoadConfig() *Config {
	return &Config{
		Server: struct {
			Address string
		}{
			Address: os.Getenv("SERVER_ADDRESS"),
		},
	}
}
