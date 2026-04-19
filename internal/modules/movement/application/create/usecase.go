package create

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/ports"
	"github.com/zchelalo/expense-control-back/internal/shared/clock"
	"github.com/zchelalo/expense-control-back/internal/shared/idgen"
	"go.uber.org/zap"
)

type UseCase struct {
	movements     ports.MovementRepository
	users         ports.UserRepository
	accounts      ports.AccountRepository
	movementTypes ports.MovementTypeRepository
	categories    ports.CategoryRepository
	clock         clock.Clock
	ids           idgen.Generator
}

func New(
	movements ports.MovementRepository,
	users ports.UserRepository,
	accounts ports.AccountRepository,
	movementTypes ports.MovementTypeRepository,
	categories ports.CategoryRepository,
	clock clock.Clock,
	ids idgen.Generator,
) *UseCase {
	return &UseCase{
		movements:     movements,
		users:         users,
		accounts:      accounts,
		movementTypes: movementTypes,
		categories:    categories,
		clock:         clock,
		ids:           ids,
	}
}

type Result struct {
	Movement *domain.Movement
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)

	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		log.Warn("invalid user ID in create movement request",
			zap.String("stage", "validate_input"),
			zap.String("user_id", cmd.UserID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	accountID, err := domain.NewAccountID(cmd.AccountID)
	if err != nil {
		log.Warn("invalid account ID in create movement request",
			zap.String("stage", "validate_input"),
			zap.String("account_id", cmd.AccountID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	movementTypeID, err := domain.NewMovementTypeID(cmd.MovementTypeID)
	if err != nil {
		log.Warn("invalid movement type ID in create movement request",
			zap.String("stage", "validate_input"),
			zap.String("movement_type_id", cmd.MovementTypeID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	categoryID, err := domain.NewCategoryID(cmd.CategoryID)
	if err != nil {
		log.Warn("invalid category ID in create movement request",
			zap.String("stage", "validate_input"),
			zap.String("category_id", cmd.CategoryID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	amount, err := domain.NewAmount(cmd.Amount)
	if err != nil {
		log.Warn("invalid amount in create movement request",
			zap.String("stage", "validate_input"),
			zap.Float64("amount", cmd.Amount),
			zap.Error(err),
		)
		return Result{}, err
	}

	description, err := domain.NewDescription(cmd.Description)
	if err != nil {
		log.Warn("invalid description in create movement request",
			zap.String("stage", "validate_input"),
			zap.String("description", cmd.Description),
			zap.Error(err),
		)
		return Result{}, err
	}

	exists, err := uc.users.Exists(ctx, userID)
	if err != nil {
		log.Error("failed to check if user exists",
			zap.String("stage", "check_user_exists"),
			zap.Error(err),
		)
		return Result{}, err
	}
	if !exists {
		log.Warn("user not found for create movement request",
			zap.String("stage", "check_user_exists"),
			zap.String("user_id", userID.String()),
		)
		return Result{}, ports.ErrNotFound{Name: "user"}
	}

	accountExists, err := uc.accounts.ExistsByUserID(ctx, accountID, userID)
	if err != nil {
		log.Error("failed to check if account belongs to user",
			zap.String("stage", "check_account_exists"),
			zap.Error(err),
		)
		return Result{}, err
	}
	if !accountExists {
		log.Warn("account not found for create movement request",
			zap.String("stage", "check_account_exists"),
			zap.String("account_id", accountID.String()),
			zap.String("user_id", userID.String()),
		)
		return Result{}, ports.ErrNotFound{Name: "account"}
	}

	movementType, err := uc.movementTypes.ByID(ctx, movementTypeID)
	if err != nil {
		log.Error("failed to get movement type",
			zap.String("stage", "get_movement_type"),
			zap.Error(err),
		)
		return Result{}, err
	}

	if _, err := uc.categories.ByID(ctx, categoryID); err != nil {
		log.Error("failed to get category",
			zap.String("stage", "get_category"),
			zap.Error(err),
		)
		return Result{}, err
	}

	now := uc.clock.Now()

	movementID, err := domain.NewMovementID(uc.ids.NewUUID())
	if err != nil {
		log.Error("failed to generate movement ID",
			zap.String("stage", "generate_movement_id"),
			zap.Error(err),
		)
		return Result{}, err
	}

	movement := domain.NewMovement(
		movementID,
		amount,
		description,
		movementTypeID,
		categoryID,
		accountID,
		userID,
		now,
	)

	operation, err := operationForMovementType(movementType)
	if err != nil {
		log.Error("failed to resolve movement type operation",
			zap.String("stage", "resolve_movement_type_operation"),
			zap.String("movement_type_key", movementType.Key()),
			zap.Error(err),
		)
		return Result{}, err
	}

	if err := uc.movements.Create(ctx, movement, operation); err != nil {
		log.Error("failed to create movement",
			zap.String("stage", "create_movement"),
			zap.Error(err),
		)
		return Result{}, err
	}

	return Result{Movement: &movement}, nil
}

func operationForMovementType(movementType domain.MovementType) (ports.BalanceOperation, error) {
	switch movementType.Key() {
	case domain.MovementTypeKeyIncome:
		return ports.BalanceOperationCredit, nil
	case domain.MovementTypeKeyExpense:
		return ports.BalanceOperationDebit, nil
	default:
		return "", domain.ErrInvalidMovementType
	}
}
