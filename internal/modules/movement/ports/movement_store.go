package ports

import (
	"context"
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
)

type BalanceOperation string

const (
	BalanceOperationCredit BalanceOperation = "credit"
	BalanceOperationDebit  BalanceOperation = "debit"
)

type MovementRepository interface {
	Create(ctx context.Context, movement domain.Movement, operation BalanceOperation) error
	ByID(ctx context.Context, movementID domain.MovementID) (domain.Movement, error)
	Delete(ctx context.Context, movement domain.Movement, operation BalanceOperation, now time.Time) error
}
