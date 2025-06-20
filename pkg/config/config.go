package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort      string
	ServerHost      string
	PostgresDSN     string
	SqliteDSN       string
	MaxHeaderBytes  int
	AllowedOrigins  []string
	Environment     string
	LogLevel        string
	RateLimit       int
	RateLimitWindow string
}

func Load() *Config {
	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "9080"),
		ServerHost:     getEnv("SERVER_HOST", "0.0.0.0"),
		PostgresDSN:    getEnv("POSTGRES_DSN", "postgres://user:password@localhost:5432/ibge?sslmode=disable"),
		SqliteDSN:      getEnv("SQLITE_DSN", "./data/ibge.db"),
		MaxHeaderBytes: getEnvAsInt("MAX_HEADER_BYTES", 1048576), // 1MB
		AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{
			"http://127.0.0.1:3000",
			"http://localhost:3000",
			"http://127.0.0.1:3001",
			"http://localhost:3001",
		}),
		Environment:     getEnv("GO_ENV", "development"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		RateLimit:       getEnvAsInt("RATE_LIMIT", 100),
		RateLimitWindow: getEnv("RATE_LIMIT_WINDOW", "1m"),
	}
}

// getEnv retorna o valor da variável de ambiente ou o valor padrão
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvAsSlice converte uma variável de ambiente em slice
func getEnvAsSlice(name string, defaultVal []string) []string {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultVal
	}
	// Remove espaços em branco e divide por vírgula
	parts := strings.Split(valStr, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// getEnvAsInt converte uma variável de ambiente em int
func getEnvAsInt(name string, defaultVal int) int {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultVal
	}
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

// GetServerAddress retorna o endereço completo do servidor
func (c *Config) GetServerAddress() string {
	return c.ServerHost + ":" + c.ServerPort
}

// IsDevelopment verifica se está em ambiente de desenvolvimento
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development" || c.Environment == "dev"
}

// IsProduction verifica se está em ambiente de produção
func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}
