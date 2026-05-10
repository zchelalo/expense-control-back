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
	movementstatsuc "github.com/zchelalo/expense-control-back/internal/modules/movement/application/stats"
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
	userReferenceStore := accountpg.NewUserReferenceRepository(db)

	createAccountUseCase := createuc.New(accountStore, userReferenceStore, clock, ids)
	listAccountsUseCase := listuc.New(accountStore, userReferenceStore, ids, cfg.PaginatorLimitDefault)
	byIDAccountsUseCase := byiduc.New(accountStore, userReferenceStore)
	updateAccountNameUseCase := updatenameuc.New(accountStore, userReferenceStore, clock)
	deleteAccountUseCase := deleteuc.New(accountStore, userReferenceStore, clock)
	accountV1 := accounthttp.NewRouter(
		createAccountUseCase,
		listAccountsUseCase,
		byIDAccountsUseCase,
		updateAccountNameUseCase,
		deleteAccountUseCase,
		mdw,
	)

	categoryStore := categorypg.NewCategoryRepo(db)
	categoryUserReferenceStore := categorypg.NewUserReferenceRepository(db)
	createCategoryUseCase := categorycreateuc.New(categoryStore, categoryUserReferenceStore, clock, ids)
	deleteCategoryUseCase := categorydeleteuc.New(categoryStore, categoryUserReferenceStore, clock)
	listCategoriesUseCase := categorylistuc.New(categoryStore, categoryUserReferenceStore, cfg.PaginatorLimitDefault)
	categoryV1 := categoryhttp.NewRouter(
		createCategoryUseCase,
		deleteCategoryUseCase,
		listCategoriesUseCase,
		mdw,
	)

	movementStore := movementpg.NewMovementRepo(db)
	movementQuery := movementpg.NewQueryRepo(db)
	movementUserReferenceStore := movementpg.NewUserReferenceRepository(db)
	movementAccountReferenceStore := movementpg.NewAccountReferenceRepository(db)
	movementTypeReferenceStore := movementpg.NewMovementTypeReferenceRepository(db)
	movementCategoryReferenceStore := movementpg.NewCategoryReferenceRepository(db)

	createMovementUseCase := movementcreateuc.New(
		movementStore,
		movementUserReferenceStore,
		movementAccountReferenceStore,
		movementTypeReferenceStore,
		movementCategoryReferenceStore,
		clock,
		ids,
	)
	listMovementsUseCase := movementlistuc.New(
		movementQuery,
		movementUserReferenceStore,
		cfg.PaginatorLimitDefault,
	)
	statsMovementsUseCase := movementstatsuc.New(
		movementQuery,
		movementUserReferenceStore,
	)
	byIDMovementUseCase := movementbyiduc.New(
		movementQuery,
		movementUserReferenceStore,
	)
	deleteMovementUseCase := movementdeleteuc.New(
		movementStore,
		movementQuery,
		movementUserReferenceStore,
		clock,
	)
	movementV1 := movementhttp.NewRouter(
		createMovementUseCase,
		listMovementsUseCase,
		byIDMovementUseCase,
		deleteMovementUseCase,
		statsMovementsUseCase,
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
