package delete

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/category/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/category/ports"
	"github.com/zchelalo/expense-control-back/internal/shared/clock"
	"go.uber.org/zap"
)

type UseCase struct {
	categories ports.CategoryRepository
	users      ports.UserRepository
	clock      clock.Clock
}

func New(
	categories ports.CategoryRepository,
	users ports.UserRepository,
	clock clock.Clock,
) *UseCase {
	return &UseCase{
		categories: categories,
		users:      users,
		clock:      clock,
	}
}

type Result struct {
	Success bool
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)

	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		log.Warn("invalid user ID in delete category request",
			zap.String("stage", "validate_input"),
			zap.String("user_id", cmd.UserID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	categoryID, err := domain.NewCategoryID(cmd.CategoryID)
	if err != nil {
		log.Warn("invalid category ID in delete category request",
			zap.String("stage", "validate_input"),
			zap.String("category_id", cmd.CategoryID.String()),
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
		log.Warn("user not found for delete category request",
			zap.String("stage", "check_user_exists"),
			zap.String("user_id", userID.String()),
		)
		return Result{}, ports.ErrNotFound{Name: "user"}
	}

	if err := uc.categories.Delete(ctx, userID, categoryID, uc.clock.Now()); err != nil {
		log.Error("failed to delete category",
			zap.String("stage", "delete_category"),
			zap.Error(err),
		)
		return Result{}, err
	}

	return Result{Success: true}, nil
}
