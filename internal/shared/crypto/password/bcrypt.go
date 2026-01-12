package password

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct {
	cost int
}

func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptPasswordHasher{cost: cost}
}

func (bph *BcryptPasswordHasher) Hash(_ context.Context, plain string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plain), bph.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (bph *BcryptPasswordHasher) Compare(_ context.Context, plain, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}