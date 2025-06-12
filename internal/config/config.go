package config

import (
    "os"
    "strconv"
)

type Config struct {
    Database DatabaseConfig
    Server   ServerConfig
    JWT      JWTConfig
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    Name     string
}

type ServerConfig struct {
    Port string
    Mode string
}

type JWTConfig struct {
    Secret     string
    Expiration int
}

func Load() *Config {
    return &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvAsInt("DB_PORT", 5432),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "postgres"),
            Name:     getEnv("DB_NAME", "healthcare_portal"),
        },
        Server: ServerConfig{
            Port: getEnv("PORT", "8080"),
            Mode: getEnv("GIN_MODE", "debug"),
        },
        JWT: JWTConfig{
            Secret:     getEnv("JWT_SECRET", "your-secret-key"),
            Expiration: getEnvAsInt("JWT_EXPIRATION", 24),
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}