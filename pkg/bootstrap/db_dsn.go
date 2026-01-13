package bootstrap

import (
	"fmt"
	"net/url"
)

func PostgresDSN(cfg Config) (string, error) {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DBUser, cfg.DBPass),
		Host:   fmt.Sprintf("%s:%d", cfg.DBHost, cfg.DBPort),
		Path:   cfg.DBName,
	}

	q := url.Values{}

	switch cfg.Environment {
	case "production":
		q.Set("sslmode", "require")
	default:
		q.Set("sslmode", "disable")
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}