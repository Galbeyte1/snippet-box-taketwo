package config

import (
	"net/http"
	"os"
	"strconv"
)

type SessionConfig struct {
	Path     string
	MaxAge   int
	HttpOnly bool
	Secure   bool
	SameSite http.SameSite
}

func NewSessionConfig() *SessionConfig {
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

	return &SessionConfig{
		Path:     "/",
		MaxAge:   sessionMaxAge,
		HttpOnly: sessionHttpOnly,
		Secure:   sessionSecure,
		SameSite: sameSite,
	}
}
