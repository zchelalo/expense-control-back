package bootstrap

import (
	"fmt"
	"os"
	"sync"

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
}

func LoadConfig(dotenvPath string) (Config, error) {
	onceConfig.Do(func() {
		loadErr = k.Load(confmap.Provider(map[string]any{
			"ENVIRONMENT": "development",
			"SERVICE_NAME": "expense-control-back",
			"PORT": 8000,
			"PAGINATOR_LIMIT_DEFAULT": 10,
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
		}
	})
	return config, loadErr
}

func GetConfig() Config { return config }