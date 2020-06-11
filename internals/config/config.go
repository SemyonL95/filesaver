package config

import "os"

type Config struct {
	AppPort     string
	RedisPort   string
	RedisHost   string
	StoragePath string
}

func NewConfig() *Config {
	return &Config{
		AppPort:     os.Getenv("APP_PORT"),
		RedisPort:   os.Getenv("REDIS_PORT"),
		RedisHost:   os.Getenv("REDIS_HOST"),
		StoragePath: os.Getenv("STORAGE_PATH"),
	}
}
