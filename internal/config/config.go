package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel                  string   `env:"LOG_LEVEL" envDefault:"info"`
	LogMode                   string   `env:"LOG_MODE" envDefault:"json"`
	Address                   string   `env:"ADDR" envDefault:":8080"`
	Secret                    string   `env:"SECRET"`
	StorageBackend            string   `env:"STORAGE_BACKEND" envDefault:"memory"`
	RedisUrl                  string   `env:"REDIS_URL"`
	HappyEmoji                string   `env:"HAPPY_EMOJI"`
	SuccessEmoji              string   `env:"SUCCESS_EMOJI"`
	FailureEmoji              string   `env:"FAILURE_EMOJI"`
	DisabledEvents            []string `env:"DISABLED_EVENTS"`
	IgnorePrivateRepositories bool     `env:"IGNORE_PRIVATE_REPOSITORIES" envDefault:"false"`
	IgnoredRepositories       []string `env:"IGNORED_REPOSITORIES"`
	IgnoredChecks             []string `env:"IGNORED_CHECKS"`
	IgnoredWorkflows          []string `env:"IGNORED_WORKFLOWS"`
}

var cfg *Config

func Init() {
	// Load env variables from file
	_ = godotenv.Load()

	cfg = &Config{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}
}

func Get() *Config {
	return cfg
}
