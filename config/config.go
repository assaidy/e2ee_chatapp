package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

var (
	ListenAddr                              = getEnvString("LISTEN_ADDR")
	AppBaseUrl                              = getEnvString("APP_BASE_URL")
	PGUrl                                   = getEnvString("PG_URL")
	EmailFrom                               = getEnvString("EMAIL_FROM")
	PapercutSmtpHost                        = getEnvString("PAPERCUT_SMTP_HOST")
	EmailVerificationTokenExpiration        = time.Hour * 24
	EmailVerificationTokenCleanupWorkerTick = time.Hour
	SessionExpiration                       time.Duration
)

func getEnvString(key string, defaultValue ...string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if len(defaultValue) == 0 {
		panic(fmt.Sprintf("env var: %s not found", key))
	}
	return defaultValue[0]
}

func getEnvInt(key string, defaultValue ...int) int {
	if value, ok := os.LookupEnv(key); ok {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			panic(fmt.Sprintf("invalid int value: %s for env var: %s", value, key))
		}
		return intValue
	}
	if len(defaultValue) == 0 {
		panic(fmt.Sprintf("env var: %s not found", key))
	}
	return defaultValue[0]
}
