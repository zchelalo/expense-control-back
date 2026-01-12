package password

import "context"

type PasswordHasher interface {
	Hash(ctx context.Context, plain string) (string, error)
	Compare(ctx context.Context, plain, hash string) error
}

