package bootstrap

import (
	"context"
	"fmt"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	authhttp "github.com/zchelalo/expense-control-back/internal/modules/auth/adapters/http/v1"
	authpg "github.com/zchelalo/expense-control-back/internal/modules/auth/adapters/persistence/postgres"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/adapters/tokens/jwt"
	loginuc "github.com/zchelalo/expense-control-back/internal/modules/auth/application/login"
	logoutuc "github.com/zchelalo/expense-control-back/internal/modules/auth/application/logout"
	refreshuc "github.com/zchelalo/expense-control-back/internal/modules/auth/application/refresh"
	registeruc "github.com/zchelalo/expense-control-back/internal/modules/auth/application/register"
	"github.com/zchelalo/expense-control-back/internal/server"
	clk "github.com/zchelalo/expense-control-back/internal/shared/clock"
	bcrypthasher "github.com/zchelalo/expense-control-back/internal/shared/crypto/password"
	uuidgenerator "github.com/zchelalo/expense-control-back/internal/shared/idgen"
	"go.uber.org/zap"
)

type App struct {
	Server  *server.Server
	Cleanup func(context.Context) error
}

func InitApp(log *zap.Logger, cfg Config) (*App, error) {
	db, err := InitDB(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot init db: %w", err)
	}

	keys, err := jwt.LoadKeys(jwt.KeyPaths{
		AccessPrivatePath: cfg.JWTAccessPrivateKeyPath,
		AccessPublicPath: cfg.JWTAccessPublicKeyPath,
		RefreshPrivatePath: cfg.JWTRefreshPrivateKeyPath,
		RefreshPublicPath: cfg.JWTRefreshPublicKeyPath,
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("cannot load jwt keys: %w", err)
	}

	clock := clk.New()
	issuer := jwt.NewIssuer(keys, clock, cfg.ServiceName, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	verifier := jwt.NewVerifier(keys, cfg.ServiceName)

	mdw := middleware.New(log, cfg.AllowedOrigins, verifier)
	address := fmt.Sprintf("0.0.0.0:%d", cfg.Port)

	hasher := bcrypthasher.NewBcryptPasswordHasher(12)

	ids := uuidgenerator.NewGenerator()

	credentialsStore := authpg.NewCredentialRepo(db)
	sessionStore := authpg.NewSessionRepo(db)

	registerUseCase := registeruc.New(credentialsStore, sessionStore, hasher, issuer, ids, clock, cfg.RefreshTokenTTL)
	loginUseCase := loginuc.New(credentialsStore, sessionStore, hasher, issuer, ids, clock, cfg.RefreshTokenTTL)
	logoutUseCase := logoutuc.New(verifier, sessionStore, clock)
	refreshUseCase := refreshuc.New(verifier, issuer, sessionStore, ids, clock, cfg.RefreshTokenTTL)

	secureCookies := cfg.Environment == "production"
	authV1 := authhttp.NewRouter(registerUseCase, loginUseCase, logoutUseCase, refreshUseCase, secureCookies, mdw)

	s, err := server.New(address, mdw, authV1.Register)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("cannot create server: %w", err)
	}

	return &App{
		Server: s,
		Cleanup: func(ctx context.Context) error {
			_ = s.Shutdown(ctx)
			db.Close()
			return nil
		},
	}, nil
}