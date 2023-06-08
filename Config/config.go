package config

import (
)

type config struct {
	Port     string `mapstructure:"PORT"`
	DBUrl    string `mapstructure:"DB_URL"`
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
}

