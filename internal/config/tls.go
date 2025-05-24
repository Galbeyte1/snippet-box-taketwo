package config

import "crypto/tls"

type TLSConfig struct {
	Config *tls.Config
}

func NewTLSConfig() *TLSConfig {
	return &TLSConfig{
		Config: &tls.Config{
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
			},
		},
	}
}
