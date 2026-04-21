package bootstrap

import (
	"context"
	"fmt"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	accounthttp "github.com/zchelalo/expense-control-back/internal/modules/account/adapters/http/v1"
	accountpg "github.com/zchelalo/expense-control-back/internal/modules/account/adapters/persistence/postgres"
	byiduc "github.com/zchelalo/expense-control-back/internal/modules/account/application/byid"
	createuc "github.com/zchelalo/expense-control-back/internal/modules/account/application/create"
	deleteuc "github.com/zchelalo/expense-control-back/internal/modules/account/application/delete"
	listuc "github.com/zchelalo/expense-control-back/internal/modules/account/application/list"
	updatenameuc "github.com/zchelalo/expense-control-back/internal/modules/account/application/updatename"
	authhttp "github.com/zchelalo/expense-control-back/internal/modules/auth/adapters/http/v1"
	authpg "github.com/zchelalo/expense-control-back/internal/modules/auth/adapters/persistence/postgres"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/adapters/tokens/jwt"
	loginuc "github.com/zchelalo/expense-control-back/internal/modules/auth/application/login"
	logoutuc "github.com/zchelalo/expense-control-back/internal/modules/auth/application/logout"
	refreshuc "github.com/zchelalo/expense-control-back/internal/modules/auth/application/refresh"
	registeruc "github.com/zchelalo/expense-control-back/internal/modules/auth/application/register"
	categoryhttp "github.com/zchelalo/expense-control-back/internal/modules/category/adapters/http/v1"
	categorypg "github.com/zchelalo/expense-control-back/internal/modules/category/adapters/persistence/postgres"
	categorycreateuc "github.com/zchelalo/expense-control-back/internal/modules/category/application/create"
	categorydeleteuc "github.com/zchelalo/expense-control-back/internal/modules/category/application/delete"
	categorylistuc "github.com/zchelalo/expense-control-back/internal/modules/category/application/list"
	movementhttp "github.com/zchelalo/expense-control-back/internal/modules/movement/adapters/http/v1"
	movementpg "github.com/zchelalo/expense-control-back/internal/modules/movement/adapters/persistence/postgres"
	movementbyiduc "github.com/zchelalo/expense-control-back/internal/modules/movement/application/byid"
	movementcreateuc "github.com/zchelalo/expense-control-back/internal/modules/movement/application/create"
	movementdeleteuc "github.com/zchelalo/expense-control-back/internal/modules/movement/application/delete"
	movementlistuc "github.com/zchelalo/expense-control-back/internal/modules/movement/application/list"
	movementtypehttp "github.com/zchelalo/expense-control-back/internal/modules/movementtype/adapters/http/v1"
	movementtypepg "github.com/zchelalo/expense-control-back/internal/modules/movementtype/adapters/persistence/postgres"
	movementtypelistuc "github.com/zchelalo/expense-control-back/internal/modules/movementtype/application/list"
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
		AccessPrivatePath:  cfg.JWTAccessPrivateKeyPath,
		AccessPublicPath:   cfg.JWTAccessPublicKeyPath,
		RefreshPrivatePath: cfg.JWTRefreshPrivateKeyPath,
		RefreshPublicPath:  cfg.JWTRefreshPublicKeyPath,
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

	accountStore := accountpg.NewAccountRepo(db)
	userStore := accountpg.NewUserRepo(db)

	createAccountUseCase := createuc.New(accountStore, userStore, clock, ids)
	listAccountsUseCase := listuc.New(accountStore, userStore, ids, cfg.PaginatorLimitDefault)
	byIDAccountsUseCase := byiduc.New(accountStore, userStore)
	updateAccountNameUseCase := updatenameuc.New(accountStore, userStore, clock)
	deleteAccountUseCase := deleteuc.New(accountStore, userStore, clock)
	accountV1 := accounthttp.NewRouter(
		createAccountUseCase,
		listAccountsUseCase,
		byIDAccountsUseCase,
		updateAccountNameUseCase,
		deleteAccountUseCase,
		mdw,
	)

	categoryStore := categorypg.NewCategoryRepo(db)
	categoryUserStore := categorypg.NewUserRepo(db)
	createCategoryUseCase := categorycreateuc.New(categoryStore, categoryUserStore, clock, ids)
	deleteCategoryUseCase := categorydeleteuc.New(categoryStore, categoryUserStore, clock)
	listCategoriesUseCase := categorylistuc.New(categoryStore, categoryUserStore, cfg.PaginatorLimitDefault)
	categoryV1 := categoryhttp.NewRouter(
		createCategoryUseCase,
		deleteCategoryUseCase,
		listCategoriesUseCase,
		mdw,
	)

	movementStore := movementpg.NewMovementRepo(db)
	movementQuery := movementpg.NewQueryRepo(db)
	movementUserStore := movementpg.NewUserRepo(db)
	movementAccountStore := movementpg.NewAccountRepo(db)
	movementTypeStore := movementpg.NewMovementTypeRepo(db)
	movementCategoryStore := movementpg.NewCategoryRepo(db)

	createMovementUseCase := movementcreateuc.New(
		movementStore,
		movementUserStore,
		movementAccountStore,
		movementTypeStore,
		movementCategoryStore,
		clock,
		ids,
	)
	listMovementsUseCase := movementlistuc.New(
		movementQuery,
		movementUserStore,
		cfg.PaginatorLimitDefault,
	)
	byIDMovementUseCase := movementbyiduc.New(
		movementQuery,
		movementUserStore,
	)
	deleteMovementUseCase := movementdeleteuc.New(
		movementStore,
		movementQuery,
		movementUserStore,
		clock,
	)
	movementV1 := movementhttp.NewRouter(
		createMovementUseCase,
		listMovementsUseCase,
		byIDMovementUseCase,
		deleteMovementUseCase,
		mdw,
	)

	movementTypesCatalogStore := movementtypepg.NewMovementTypeRepo(db)
	listMovementTypesUseCase := movementtypelistuc.New(movementTypesCatalogStore)
	movementTypeV1 := movementtypehttp.NewRouter(
		listMovementTypesUseCase,
		mdw,
	)

	s, err := server.New(address, mdw, authV1.Register, accountV1.Register, categoryV1.Register, movementV1.Register, movementTypeV1.Register)
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
