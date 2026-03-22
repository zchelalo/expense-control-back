package create

import "github.com/zchelalo/expense-control-back/internal/modules/account/domain"

type Command struct {
	Name   domain.Name
	Balance domain.Balance
	UserID domain.UserID
}