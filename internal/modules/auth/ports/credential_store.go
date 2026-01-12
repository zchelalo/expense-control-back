package ports

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/modules/auth/domain"
)

type CredentialStore interface {
	CreateAccount(ctx context.Context, email domain.Email, passHash domain.PasswordHash) (domain.SubjectID, error)
	ByEmail(ctx context.Context, email domain.Email) (domain.Account, error)
}