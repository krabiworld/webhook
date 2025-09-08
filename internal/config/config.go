package config

import "os"

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
	cfg = &Config{
		LogLevel:     env("LOG_LEVEL", "info"),
		Address:      env("ADDR", ":8080"),
		Secret:       env("SECRET"),
		HappyEmoji:   env("HAPPY_EMOJI"),
		SuccessEmoji: env("SUCCESS_EMOJI"),
		FailureEmoji: env("FAILURE_EMOJI"),
	}
}

func env(key string, defaultValue ...string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func Get() *Config {
	return cfg
}
