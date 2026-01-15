package logout

import "github.com/zchelalo/expense-control-back/internal/modules/auth/domain"

type Command struct {
	SubjectID domain.SubjectID
	RefreshToken string
}