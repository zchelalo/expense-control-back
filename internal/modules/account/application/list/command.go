package list

import (
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
)

type Command struct {
	UserID    domain.UserID
	CreatedAt *time.Time
	AccountID *domain.AccountID
	Limit     int
	IsBefore  bool
}