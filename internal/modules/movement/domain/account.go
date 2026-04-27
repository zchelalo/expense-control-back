package domain

type Account struct {
	id   AccountID
	name string
}

func RehydrateAccount(id AccountID, name string) (Account, error) {
	if name == "" {
		return Account{}, ErrInvalidAccount
	}

	return Account{
		id:   id,
		name: name,
	}, nil
}

func (a Account) ID() AccountID { return a.id }
func (a Account) Name() string  { return a.name }
