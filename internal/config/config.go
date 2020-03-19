package config

import (
	"github.com/kelseyhightower/envconfig"
	"os"
)

const envFileName = ".env"
const devEnv = "dev"

func init() {
	if os.Getenv("SERVICE_ENV") == devEnv {
		MustLoadEnv()
	}
}

// NewConfig returns the settings from the environment.
func NewConfig() *Config {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func (c Config) IsDev() bool {
	return c.Service.Env == devEnv
}
