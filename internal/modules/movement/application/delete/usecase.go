package delete

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/ports"
	"github.com/zchelalo/expense-control-back/internal/shared/clock"
	"go.uber.org/zap"
)

type UseCase struct {
	movements ports.MovementRepository
	query     ports.QueryRepository
	users     ports.UserRepository
	clock     clock.Clock
}

func New(
	movements ports.MovementRepository,
	query ports.QueryRepository,
	users ports.UserRepository,
	clock clock.Clock,
) *UseCase {
	return &UseCase{
		movements: movements,
		query:     query,
		users:     users,
		clock:     clock,
	}
}

type Result struct {
	Success bool
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)

	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		log.Warn("invalid user ID in delete movement request",
			zap.String("stage", "validate_input"),
			zap.String("user_id", cmd.UserID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	movementID, err := domain.NewMovementID(cmd.MovementID)
	if err != nil {
		log.Warn("invalid movement ID in delete movement request",
			zap.String("stage", "validate_input"),
			zap.String("movement_id", cmd.MovementID.String()),
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
		log.Warn("user not found for delete movement request",
			zap.String("stage", "check_user_exists"),
			zap.String("user_id", userID.String()),
		)
		return Result{}, ports.ErrNotFound{Name: "user"}
	}

	movementDetails, err := uc.query.ByIDForUser(ctx, movementID, userID)
	if err != nil {
		log.Error("failed to get movement by id",
			zap.String("stage", "get_movement_by_id"),
			zap.Error(err),
		)
		return Result{}, err
	}

	operation, err := reverseOperationForMovementType(movementDetails.MovementType())
	if err != nil {
		log.Error("failed to resolve reverse operation for delete movement",
			zap.String("stage", "resolve_reverse_operation"),
			zap.String("movement_type_key", movementDetails.MovementType().Key()),
			zap.Error(err),
		)
		return Result{}, err
	}

	now := uc.clock.Now()
	if err := uc.movements.Delete(ctx, movementDetails.Movement(), operation, now); err != nil {
		log.Error("failed to delete movement",
			zap.String("stage", "delete_movement"),
			zap.Error(err),
		)
		return Result{}, err
	}

	return Result{Success: true}, nil
}

func reverseOperationForMovementType(movementType domain.MovementType) (ports.BalanceOperation, error) {
	switch movementType.Key() {
	case domain.MovementTypeKeyIncome:
		return ports.BalanceOperationDebit, nil
	case domain.MovementTypeKeyExpense:
		return ports.BalanceOperationCredit, nil
	default:
		return "", domain.ErrInvalidMovementType
	}
}
