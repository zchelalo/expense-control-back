package updatename

import "github.com/zchelalo/expense-control-back/internal/modules/account/domain"

type Command struct {
	UserID    domain.UserID
	AccountID domain.AccountID
	Name 		  domain.Name
}