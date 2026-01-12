package domain

type PasswordHash struct{ value string }

func NewPasswordHashFromHash(hash string) (PasswordHash, error) {
	if hash == "" {
		return PasswordHash{}, ErrInvalidPasswordHash
	}
	return PasswordHash{value: hash}, nil
}

func (p PasswordHash) String() string { return p.value }