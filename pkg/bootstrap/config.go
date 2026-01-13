package bootstrap

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

var (
	config     Config
	onceConfig sync.Once
	loadErr    error
)

type Config struct {
	Environment 										  string
	ServiceName											  string
	Port                              int

	DBHost                            string
	DBUser                            string
	DBPass                            string
	DBName                            string
	DBPort                            int

	PaginatorLimitDefault             int

	AllowedOrigins                    string

	OtelExporterOtlpEndpoint 				  string

	AccessTokenTTL  									time.Duration
	RefreshTokenTTL 									time.Duration

	JWTAccessPrivateKeyPath   				string
	JWTAccessPublicKeyPath    				string
	JWTRefreshPrivateKeyPath  				string
	JWTRefreshPublicKeyPath   				string
}

func LoadConfig(dotenvPath string) (Config, error) {
	onceConfig.Do(func() {
		loadErr = k.Load(confmap.Provider(map[string]any{
			"ENVIRONMENT": "development",
			"SERVICE_NAME": "expense-control-back",
			"PORT": 8000,

			"PAGINATOR_LIMIT_DEFAULT": 10,

			"OTEL_EXPORTER_OTLP_ENDPOINT": "expense-control-otel-collector:4317",

			"ACCESS_TOKEN_TTL":  "15m",
			"REFRESH_TOKEN_TTL": "720h",
		}, "."), nil)
		if loadErr != nil { return }

		if dotenvPath != "" {
			if _, err := os.Stat(dotenvPath); err == nil {
				loadErr = k.Load(file.Provider(dotenvPath), dotenv.Parser())
				if loadErr != nil {
					loadErr = fmt.Errorf("error loading %s: %w", dotenvPath, loadErr)
					return
				}
			} else if !os.IsNotExist(err) {
				loadErr = fmt.Errorf("error checking %s: %w", dotenvPath, err)
				return
			}
		}

		loadErr = k.Load(env.Provider("", ".", nil), nil)
		if loadErr != nil {
			return
		}

		shouldExistKeys := []string{
			"DB_HOST",
			"DB_USER",
			"DB_PASS",
			"DB_NAME",
			"DB_PORT",
			"ALLOWED_ORIGINS",
			"JWT_ACCESS_PRIVATE_KEY_PATH",
			"JWT_ACCESS_PUBLIC_KEY_PATH",
			"JWT_REFRESH_PRIVATE_KEY_PATH",
			"JWT_REFRESH_PUBLIC_KEY_PATH",
		}

		for _, key := range shouldExistKeys {
			if !k.Exists(key) {
				loadErr = fmt.Errorf("missing required config key: %s", key)
				return
			}
		}

		port := k.Int("PORT")
		if port <= 0 {
			loadErr = fmt.Errorf("PORT must be a valid number > 0")
			return
		}

		dbPort := k.Int("DB_PORT")
		if dbPort <= 0 {
			loadErr = fmt.Errorf("DB_PORT must be a valid number > 0")
			return
		}

		paginatorLimitDefault := k.Int("PAGINATOR_LIMIT_DEFAULT")
		if paginatorLimitDefault <= 0 {
			loadErr = fmt.Errorf("PAGINATOR_LIMIT_DEFAULT must be a valid number > 0")
			return
		}

		accessTTLStr := k.String("ACCESS_TOKEN_TTL")
		accessTTL, err := time.ParseDuration(accessTTLStr)
		if err != nil || accessTTL <= 0 {
			loadErr = fmt.Errorf("ACCESS_TOKEN_TTL must be a valid duration (e.g. 15m)")
			return
		}

		refreshTTLStr := k.String("REFRESH_TOKEN_TTL")
		refreshTTL, err := time.ParseDuration(refreshTTLStr)
		if err != nil || refreshTTL <= 0 {
			loadErr = fmt.Errorf("REFRESH_TOKEN_TTL must be a valid duration (e.g. 720h)")
			return
		}

		if accessTTL >= refreshTTL {
			loadErr = fmt.Errorf("ACCESS_TOKEN_TTL must be smaller than REFRESH_TOKEN_TTL")
			return
		}

		if accessTTL > time.Hour {
			loadErr = fmt.Errorf("ACCESS_TOKEN_TTL too large (recommended <= 1h)")
			return
		}

		accessPrivateKeyPath := k.String("JWT_ACCESS_PRIVATE_KEY_PATH")
		if _, err := os.Stat(accessPrivateKeyPath); err != nil {
			loadErr = fmt.Errorf("JWT_ACCESS_PRIVATE_KEY_PATH invalid: %w", err)
		}

		accessPublicKeyPath := k.String("JWT_ACCESS_PUBLIC_KEY_PATH")
		if _, err := os.Stat(accessPublicKeyPath); err != nil {
			loadErr = fmt.Errorf("JWT_ACCESS_PUBLIC_KEY_PATH invalid: %w", err)
		}

		refreshPrivateKeyPath := k.String("JWT_REFRESH_PRIVATE_KEY_PATH")
		if _, err := os.Stat(refreshPrivateKeyPath); err != nil {
			loadErr = fmt.Errorf("JWT_REFRESH_PRIVATE_KEY_PATH invalid: %w", err)
		}

		refreshPublicKeyPath := k.String("JWT_REFRESH_PUBLIC_KEY_PATH")
		if _, err := os.Stat(refreshPublicKeyPath); err != nil {
			loadErr = fmt.Errorf("JWT_REFRESH_PUBLIC_KEY_PATH invalid: %w", err)
		}

		config = Config{
			Environment:           k.String("ENVIRONMENT"),
			ServiceName:           k.String("SERVICE_NAME"),
			Port:                  port,
			DBHost:                k.String("DB_HOST"),
			DBUser:                k.String("DB_USER"),
			DBPass:                k.String("DB_PASS"),
			DBName:                k.String("DB_NAME"),
			DBPort:                dbPort,
			PaginatorLimitDefault: paginatorLimitDefault,
			AllowedOrigins:        k.String("ALLOWED_ORIGINS"),
			OtelExporterOtlpEndpoint: k.String("OTEL_EXPORTER_OTLP_ENDPOINT"),
			AccessTokenTTL: accessTTL,
			RefreshTokenTTL: refreshTTL,
			JWTAccessPrivateKeyPath:    accessPrivateKeyPath,
			JWTAccessPublicKeyPath:     accessPublicKeyPath,
			JWTRefreshPrivateKeyPath:   refreshPrivateKeyPath,
			JWTRefreshPublicKeyPath:    refreshPublicKeyPath,
		}
	})
	return config, loadErr
}

func GetConfig() Config { return config }