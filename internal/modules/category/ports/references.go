package ports

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/modules/category/domain"
)

type UserRepository interface {
	Exists(ctx context.Context, userID domain.UserID) (bool, error)
}
