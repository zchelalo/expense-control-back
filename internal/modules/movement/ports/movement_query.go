package ports

import (
	"context"
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
)

type ListMovementsFilter struct {
	AccountID      *domain.AccountID
	CategoryID     *domain.CategoryID
	MovementTypeID *domain.MovementTypeID
	CreatedAt      *time.Time
	MovementID     *domain.MovementID
	Limit          int
	IsBefore       bool
}

// QueryRepository is the read side to hydrate details.
type QueryRepository interface {
	ByIDForUser(ctx context.Context, movementID domain.MovementID, userID domain.UserID) (domain.MovementDetails, error)
	ListByUserID(ctx context.Context, userID domain.UserID, filter ListMovementsFilter) ([]domain.MovementDetails, error)
}
