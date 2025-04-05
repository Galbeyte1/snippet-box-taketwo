package config

import "log/slog"

type Application struct {
	Logger *slog.Logger
}

type Config struct {
	Addr      string
	StaticDir string
	DSN       string
}
