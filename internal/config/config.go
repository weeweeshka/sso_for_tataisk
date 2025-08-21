package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	StoragePath string        `yaml:"storage_path"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-default:"1h"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := "config/local.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not found: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic(err)
	}

	return &cfg

}
