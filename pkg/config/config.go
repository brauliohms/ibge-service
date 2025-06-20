package config

import "os"

type Config struct {
	ServerPort     string
	PostgresDSN    string
	SqliteDSN      string
	MaxHeaderBytes int
}

func Load() *Config {
	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "9080"),
		PostgresDSN:    getEnv("POSTGRES_DSN", "postgres://user:password@localhost:5432/ibge?sslmode=disable"),
		SqliteDSN:      getEnv("SQLITE_DSN", "./data/ibge.db"),
		MaxHeaderBytes: 1 << 20, // 1MB
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
