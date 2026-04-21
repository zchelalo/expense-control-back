package ports

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
)

type UserRepository interface {
	Exists(ctx context.Context, userID domain.UserID) (bool, error)
}

type AccountRepository interface {
	ExistsByUserID(ctx context.Context, accountID domain.AccountID, userID domain.UserID) (bool, error)
}

type MovementTypeRepository interface {
	ByID(ctx context.Context, movementTypeID domain.MovementTypeID) (domain.MovementType, error)
}

type CategoryRepository interface {
	ByIDForUser(ctx context.Context, categoryID domain.CategoryID, userID domain.UserID) (domain.Category, error)
}
