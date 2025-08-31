package config

import (
	"os"
	"time"
)

type jwtCfg struct {
	Secret    string
	AccessTTL time.Duration
}

var JWT = jwtCfg{
	Secret:    getenv("JWT_SECRET", "dev-secret-change-me"),
	AccessTTL: mustParseDuration(getenv("JWT_ACCESS_TTL", "15m")),
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func mustParseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}
