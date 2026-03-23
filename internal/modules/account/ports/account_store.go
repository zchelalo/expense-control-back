package ports

import (
	"context"
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
)

type AccountRepository interface {
	Create(ctx context.Context, account domain.Account) error
	ByID(ctx context.Context, accountID domain.AccountID) (domain.Account, error)
	ListByUserID(ctx context.Context, userID domain.UserID, name *string, createdAt *time.Time, accountID *domain.AccountID, limit int, isBefore bool) ([]domain.Account, error)
	UpdateName(ctx context.Context, accountID domain.AccountID, name domain.Name, now time.Time) error
	UpdateBalance(ctx context.Context, accountID domain.AccountID, balance domain.Balance, now time.Time) error
	Delete(ctx context.Context, accountID domain.AccountID, now time.Time) error
}