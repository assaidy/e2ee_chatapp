package env

import (
	"log/slog"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

var (
	ServerPort   = GetEnvInt("SERVER_PORT", 8080)
	ServerSecret = GetEnvString("SERVER_SECRET")
	PgUrl        = GetEnvString("PG_URL")
	ValkeyPort   = GetEnvInt("VALKEY_PORT", 6379)
)

func GetEnvInt(key string, defaultValue ...int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		slog.Error("invalid int env var", "key", key, "value", value)
		os.Exit(1)
	}
	if len(defaultValue) == 0 {
		slog.Error("env var not found", "key", key)
		os.Exit(1)
	}
	return defaultValue[0]
}

func GetEnvString(key string, defaultValue ...string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if len(defaultValue) == 0 {
		slog.Error("env var not found", "key", key)
		os.Exit(1)
	}
	return defaultValue[0]
}
