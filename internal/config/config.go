package config

import (
	"os"
	"time"
)

type jwtCfg struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

// JWT konfigürasyonu
func GetJWT() jwtCfg {
	return jwtCfg{
		Secret:     getenv("JWT_SECRET", "dev-secret-change-me"),
		AccessTTL:  mustParseDuration(getenv("JWT_ACCESS_TTL", "24h")),
		RefreshTTL: mustParseDuration(getenv("JWT_REFRESH_TTL", "168h")),
	}
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
		return 24 * time.Hour
	}
	return d
}
