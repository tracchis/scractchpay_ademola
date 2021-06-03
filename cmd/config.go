package main

import (
	"time"
)

// Config representation of the service configuration
type Config struct {
	*GlobalConfig
	Env  string `envconfig:"ENVIRONMENT" default:"development"`
	Port int    `envconfig:"PORT" default:"8000"`
}

// GlobalConfig represents common application parameters
type GlobalConfig struct {
	Port              int           `envconfig:"PORT"`
	ClientTimeout     int           `envconfig:"CLIENT_TIMEOUT_SEC"`
	ClientIdleTimeout time.Duration `envconfig:"CLIENT_IDLE_TIMEOUT"`
	LogLevel          string        `envconfig:"LOG_LEVEL"`
	AppEnv            string        `envconfig:"APP_ENV"`
}
