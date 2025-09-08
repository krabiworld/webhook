package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel     string `env:"LOG_LEVEL" envDefault:"info"`
	Address      string `env:"ADDR" envDefault:":8080"`
	Secret       string `env:"SECRET"`
	HappyEmoji   string `env:"HAPPY_EMOJI"`
	SuccessEmoji string `env:"SUCCESS_EMOJI"`
	FailureEmoji string `env:"FAILURE_EMOJI"`
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
