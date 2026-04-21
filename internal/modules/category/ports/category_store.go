package ports

import (
	"context"
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/category/domain"
)

type CategoryRepository interface {
	Create(ctx context.Context, category domain.Category) (domain.Category, error)
	ListByUserID(ctx context.Context, userID domain.UserID, createdAt *time.Time, categoryID *domain.CategoryID, limit int, isBefore bool) ([]domain.Category, error)
	Delete(ctx context.Context, userID domain.UserID, categoryID domain.CategoryID, now time.Time) error
}
