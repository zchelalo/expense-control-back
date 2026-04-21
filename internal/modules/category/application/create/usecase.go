package create

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/category/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/category/ports"
	"github.com/zchelalo/expense-control-back/internal/shared/clock"
	"github.com/zchelalo/expense-control-back/internal/shared/idgen"
	"go.uber.org/zap"
)

type UseCase struct {
	categories ports.CategoryRepository
	users      ports.UserRepository
	clock      clock.Clock
	ids        idgen.Generator
}

func New(
	categories ports.CategoryRepository,
	users ports.UserRepository,
	clock clock.Clock,
	ids idgen.Generator,
) *UseCase {
	return &UseCase{
		categories: categories,
		users:      users,
		clock:      clock,
		ids:        ids,
	}
}

type Result struct {
	Category *domain.Category
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)

	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		log.Warn("invalid user ID in create category request",
			zap.String("stage", "validate_input"),
			zap.String("user_id", cmd.UserID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	name, err := domain.NewName(cmd.Name)
	if err != nil {
		log.Warn("invalid category name in create category request",
			zap.String("stage", "validate_input"),
			zap.String("name", cmd.Name),
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
		log.Warn("user not found for create category request",
			zap.String("stage", "check_user_exists"),
			zap.String("user_id", userID.String()),
		)
		return Result{}, ports.ErrNotFound{Name: "user"}
	}

	categoryID, err := domain.NewCategoryID(uc.ids.NewUUID())
	if err != nil {
		log.Error("failed to generate category ID",
			zap.String("stage", "generate_category_id"),
			zap.Error(err),
		)
		return Result{}, err
	}

	category := domain.NewCategory(
		categoryID,
		name,
		userID,
		uc.clock.Now(),
	)

	createdCategory, err := uc.categories.Create(ctx, category)
	if err != nil {
		log.Error("failed to create category",
			zap.String("stage", "create_category"),
			zap.Error(err),
		)
		return Result{}, err
	}

	return Result{Category: &createdCategory}, nil
}
