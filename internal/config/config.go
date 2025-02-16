// Package config Description: Определение структуры конфигурации приложения
// и метода для загрузки конфигурации из переменных окружения.
package config

import "os"

// Config настройка приложения
type Config struct {
	Server struct {
		Address string
	}
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	return &Config{
		Server: struct {
			Address string
		}{
			Address: os.Getenv("SERVER_ADDRESS"),
		},
	}
}
