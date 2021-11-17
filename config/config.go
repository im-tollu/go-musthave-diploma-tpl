// Package config is responsible for taking the runtime configuration from
// multiple sources of parameters and providing a structured configuration
// data to the service at the time of launch. It is also provides sensible
// defaults.
//
// Environment variables are considered the primary source of configuration.
// It supports the 12-factors app approach.
// For developers' convenience configuration can be overridden
// with CLI parameters.
package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:"localhost:8081"`
	DatabaseURI          string `env:"DATABASE_URI" envDefault:"postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:8080"`
}

func Load() (*Config, error) {
	config := &Config{}

	if errEnv := env.Parse(config); errEnv != nil {
		return nil, fmt.Errorf("cannot parse config from environment: %w", errEnv)
	}

	overrideWithCliParams(config)

	return config, nil
}

func overrideWithCliParams(config *Config) {
	flag.StringVar(&config.RunAddress, "a", config.RunAddress, "адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a")
	flag.StringVar(&config.DatabaseURI, "d", config.DatabaseURI, "адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d")
	flag.StringVar(&config.AccrualSystemAddress, "r", config.AccrualSystemAddress, "адрес системы расчёта начислений: переменная окружения ОС ACCRUAL_SYSTEM_ADDRESS или флаг -r")

	flag.Parse()
}
