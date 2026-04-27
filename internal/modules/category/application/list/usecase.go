package list

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/category/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/category/ports"
	"go.uber.org/zap"
)

type UseCase struct {
	categories      ports.CategoryRepository
	users           ports.UserReferenceRepository
	pagLimitDefault int
}

func New(
	categories ports.CategoryRepository,
	users ports.UserReferenceRepository,
	pagLimitDefault int,
) *UseCase {
	return &UseCase{
		categories:      categories,
		users:           users,
		pagLimitDefault: pagLimitDefault,
	}
}

type Result struct {
	Categories []domain.Category
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)

	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		log.Warn("invalid user ID in list categories request",
			zap.String("stage", "validate_input"),
			zap.String("user_id", cmd.UserID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	var categoryID *domain.CategoryID
	if cmd.CategoryID != nil {
		id, err := domain.NewCategoryID(*cmd.CategoryID)
		if err != nil {
			log.Warn("invalid category ID in list categories request",
				zap.String("stage", "validate_input"),
				zap.String("category_id", cmd.CategoryID.String()),
				zap.Error(err),
			)
			return Result{}, err
		}
		categoryID = &id
	}

	if cmd.Limit <= 0 {
		log.Warn("invalid pagination limit in list categories request",
			zap.String("stage", "validate_input"),
			zap.Int("limit", cmd.Limit),
		)
		cmd.Limit = uc.pagLimitDefault
	}

	if cmd.CreatedAt == nil && categoryID != nil {
		log.Warn("category ID provided without created at in list categories request",
			zap.String("stage", "validate_input"),
			zap.String("category_id", categoryID.String()),
		)
		return Result{}, ErrCategoryIDWithoutCreatedAt
	}

	if cmd.CreatedAt != nil && categoryID == nil {
		log.Warn("created at provided without category ID in list categories request",
			zap.String("stage", "validate_input"),
			zap.Time("created_at", *cmd.CreatedAt),
		)
		return Result{}, ErrCreatedAtWithoutCategoryID
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
		log.Warn("user not found for list categories request",
			zap.String("stage", "check_user_exists"),
			zap.String("user_id", userID.String()),
		)
		return Result{}, ports.ErrNotFound{Name: "user"}
	}

	categories, err := uc.categories.ListByUserID(ctx, userID, cmd.Name, cmd.CreatedAt, categoryID, cmd.Limit, cmd.IsBefore)
	if err != nil {
		log.Error("failed to list categories",
			zap.String("stage", "list_categories"),
			zap.Error(err),
		)
		return Result{}, err
	}

	return Result{Categories: categories}, nil
}
