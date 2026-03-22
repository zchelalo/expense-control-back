package ports

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
)

type UserRepository interface {
	Exists(ctx context.Context, userID domain.UserID) (bool, error)
}