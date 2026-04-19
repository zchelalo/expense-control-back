package byid

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/ports"
	"go.uber.org/zap"
)

type UseCase struct {
	query ports.QueryRepository
	users ports.UserRepository
}

func New(
	query ports.QueryRepository,
	users ports.UserRepository,
) *UseCase {
	return &UseCase{
		query: query,
		users: users,
	}
}

type Result struct {
	Movement domain.MovementDetails
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)

	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		log.Warn("invalid user ID in get movement by id request",
			zap.String("stage", "validate_input"),
			zap.String("user_id", cmd.UserID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	movementID, err := domain.NewMovementID(cmd.MovementID)
	if err != nil {
		log.Warn("invalid movement ID in get movement by id request",
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
		log.Warn("user not found for get movement by id request",
			zap.String("stage", "check_user_exists"),
			zap.String("user_id", userID.String()),
		)
		return Result{}, ports.ErrNotFound{Name: "user"}
	}

	movement, err := uc.query.ByIDForUser(ctx, movementID, userID)
	if err != nil {
		log.Error("failed to get movement by id",
			zap.String("stage", "get_movement_by_id"),
			zap.Error(err),
		)
		return Result{}, err
	}

	return Result{Movement: movement}, nil
}
