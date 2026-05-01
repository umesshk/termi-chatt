package config

import (
	"os"
)

type Config struct {
	Port string

	PostgresDSN string

	RedisAddr     string
	RedisPassword string
	RedisDB       string
}

func FromEnv() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Prefer a full DSN so users can run Postgres however they want.
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		// Back-compat with the repo defaults (but no longer hardcode credentials in code).
		dsn = "host=localhost port=5432 user=postgres password=mypass dbname=termichatt sslmode=disable"
	}

	return Config{
		Port:          port,
		PostgresDSN:   dsn,
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       os.Getenv("REDIS_DB"),
	}
}

