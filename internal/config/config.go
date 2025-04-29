package config

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Galbeyte1/snippet-box-taketwo/internal/database"
)

type Config struct {
	APIAddr     string
	StaticDir   string
	Database    database.DBConfig
	Redis       database.RedisConfig
	SessionKey  string
	SessionOpts SessionConfig
	Env         string
}

type SessionConfig struct {
	Path     string
	MaxAge   int
	HttpOnly bool
	Secure   bool
	SameSite http.SameSite
}

func LoadConfigFromEnv() Config {
	// Parse Redis timeouts
	idleTimeoutSeconds := parseEnvInt("REDIS_IDLE_TIMEOUT", 240)
	redisIdleTimeout := time.Duration(idleTimeoutSeconds) * time.Second

	// Parse Session cookie settings
	sessionMaxAge, _ := strconv.Atoi(os.Getenv("SESSION_MAX_AGE"))
	sessionSecure := os.Getenv("SESSION_SECURE") == "true"
	sessionHttpOnly := os.Getenv("SESSION_HTTP_ONLY") == "true"

	// Parse Session SameSite option
	var sameSite http.SameSite
	switch os.Getenv("SESSION_SAME_SITE") {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "lax":
		sameSite = http.SameSiteLaxMode
	case "none":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteLaxMode
	}

	env := os.Getenv("ENV")

	// Build config
	cfg := Config{
		APIAddr:    os.Getenv("API_ADDR"),
		StaticDir:  os.Getenv("STATIC_DIR"),
		SessionKey: os.Getenv("SESSION_SECRET"),
		Env:        env,
		Database: database.DBConfig{
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_DATABASE"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		Redis: database.RedisConfig{
			Addr:        os.Getenv("REDIS_ADDR"),
			MaxIdle:     parseEnvInt("REDIS_MAX_IDLE", 10),
			MaxActive:   parseEnvInt("REDIS_MAX_ACTIVE", 100),
			IdleTimeout: redisIdleTimeout,
		},
		SessionOpts: SessionConfig{
			Path:     "/",
			MaxAge:   sessionMaxAge,
			HttpOnly: sessionHttpOnly,
			Secure:   sessionSecure,
			SameSite: sameSite,
		},
	}

	if cfg.IsProduction() {
		cfg.SessionOpts.Secure = true
	}

	return cfg
}

func parseEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return i
}

func (c Config) IsProduction() bool {
	return c.Env == "production"
}
