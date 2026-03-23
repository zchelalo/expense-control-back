package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	accountdb "github.com/zchelalo/expense-control-back/internal/db/sqlc/account"
	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/account/ports"
)

type AccountRepo struct {
  q *accountdb.Queries
}

func NewAccountRepo(db accountdb.DBTX) *AccountRepo {
  return &AccountRepo{q: accountdb.New(db)}
}

func (r *AccountRepo) Create(ctx context.Context, s domain.Account) error {
	var deletedAt pgtype.Timestamptz
	if s.DeletedAt() != nil {
		deletedAt = pgtype.Timestamptz{Time: *s.DeletedAt(), Valid: true}
	} else {
		deletedAt = pgtype.Timestamptz{Valid: false}
	}

	var balance pgtype.Numeric
	if err := balance.Scan(fmt.Sprintf("%f", s.Balance().Float64())); err != nil {
		return err
	}

	params := accountdb.CreateAccountParams{
		ID: pgtype.UUID{
			Bytes: s.ID().UUID(),
			Valid: true,
		},
		Name: s.Name().String(),
		Balance: balance,
		UserID: pgtype.UUID{
			Bytes: s.UserID().UUID(),
			Valid: true,
		},
		CreatedAt: pgtype.Timestamptz{
			Time: s.CreatedAt(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time: s.UpdatedAt(),
			Valid: true,
		},
		DeletedAt: deletedAt,
	}

	err := r.q.CreateAccount(ctx, params)
	return err
}

func (r *AccountRepo) ByID(ctx context.Context, id domain.AccountID) (domain.Account, error) {
	account, err := r.q.GetAccountByID(ctx, pgtype.UUID{
		Bytes: id.UUID(),
		Valid: true,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Account{}, ports.ErrNotFound{Name:"account"}
		}
		return domain.Account{}, err
	}

	parsedID, err := domain.NewAccountID(account.ID.Bytes)
	if err != nil {
		return domain.Account{}, err
	}

	parsedName, err := domain.NewName(account.Name)
	if err != nil {
		return domain.Account{}, err
	}

	balance, err := account.Balance.Float64Value()
	if err != nil {
		return domain.Account{}, err
	}

	parsedBalance, err := domain.NewBalance(balance.Float64)
	if err != nil {
		return domain.Account{}, err
	}

	var parsedDeletedAt *time.Time
	if account.DeletedAt.Valid {
		t := account.DeletedAt.Time
		parsedDeletedAt = &t
	}

	parsedUserID, err := domain.NewUserID(account.UserID.Bytes)
	if err != nil {
		return domain.Account{}, err
	}

	return domain.RehydrateAccount(
		parsedID,
		parsedName,
		parsedBalance,
		parsedUserID,
		account.CreatedAt.Time,
		account.UpdatedAt.Time,
		parsedDeletedAt,
	), nil
}
func (r *AccountRepo) ListByUserID(ctx context.Context, userID domain.UserID, createdAt *time.Time, accountID *domain.AccountID, limit int, isBefore bool) ([]domain.Account, error) {
	var accounts []accountdb.Account
	var err error

	if isBefore {
		accounts, err = r.q.ListAccountsByUserIDBefore(ctx, accountdb.ListAccountsByUserIDBeforeParams{
			UserID: pgtype.UUID{
				Bytes: userID.UUID(),
				Valid: true,
			},
			Column2: pgtype.Timestamptz{
				Time:  *createdAt,
				Valid: true,
			},
			Column3: pgtype.UUID{
				Bytes: accountID.UUID(),
				Valid: true,
			},
			Limit: int32(limit),
		})
	} else {
		params := accountdb.ListAccountsByUserIDAfterParams{
			UserID: pgtype.UUID{
				Bytes: userID.UUID(),
				Valid: true,
			},
			Limit: int32(limit),
		}

		if createdAt != nil {
			params.Column2 = pgtype.Timestamptz{
				Time:  *createdAt,
				Valid: true,
			}
		}

		if accountID != nil {
			params.Column3 = pgtype.UUID{
				Bytes: accountID.UUID(),
				Valid: true,
			}
		}

		accounts, err = r.q.ListAccountsByUserIDAfter(ctx, params)
	}

	if err != nil {
		return nil, err
	}

	if isBefore {
		for i, j := 0, len(accounts)-1; i < j; i, j = i+1, j-1 {
			accounts[i], accounts[j] = accounts[j], accounts[i]
		}
	}

	result := make([]domain.Account, len(accounts))
	for i, account := range accounts {
		parsedID, err := domain.NewAccountID(account.ID.Bytes)
		if err != nil {
			return nil, err
		}

		parsedName, err := domain.NewName(account.Name)
		if err != nil {
			return nil, err
		}

		balance, err := account.Balance.Float64Value()
		if err != nil {
			return nil, err
		}

		parsedBalance, err := domain.NewBalance(balance.Float64)
		if err != nil {
			return nil, err
		}

		var parsedDeletedAt *time.Time
		if account.DeletedAt.Valid {
			t := account.DeletedAt.Time
			parsedDeletedAt = &t
		}

		parsedUserID, err := domain.NewUserID(account.UserID.Bytes)
		if err != nil {
			return nil, err
		}

		result[i] = domain.RehydrateAccount(
			parsedID,
			parsedName,
			parsedBalance,
			parsedUserID,
			account.CreatedAt.Time,
			account.UpdatedAt.Time,
			parsedDeletedAt,
		)
	}

	return result, nil
}

func (r *AccountRepo) UpdateName(ctx context.Context, accountID domain.AccountID, name domain.Name, now time.Time) error {
	return r.q.UpdateAccountName(ctx, accountdb.UpdateAccountNameParams{
		ID: pgtype.UUID{
			Bytes: accountID.UUID(),
			Valid: true,
		},
		Name: name.String(),
		UpdatedAt: pgtype.Timestamptz{
			Time: now,
			Valid: true,
		},
	})
}

func (r *AccountRepo) UpdateBalance(ctx context.Context, accountID domain.AccountID, balance domain.Balance, now time.Time) error {
	var balanceNumeric pgtype.Numeric
	if err := balanceNumeric.Scan(fmt.Sprintf("%f", balance.Float64())); err != nil {
		return err
	}

	return r.q.UpdateAccountBalance(ctx, accountdb.UpdateAccountBalanceParams{
		ID: pgtype.UUID{
			Bytes: accountID.UUID(),
			Valid: true,
		},
		Balance: balanceNumeric,
		UpdatedAt: pgtype.Timestamptz{
			Time: now,
			Valid: true,
		},
	})
}

func (r *AccountRepo) Delete(ctx context.Context, accountID domain.AccountID, now time.Time) error {
	return r.q.DeleteAccount(ctx, accountdb.DeleteAccountParams{
		ID: pgtype.UUID{
			Bytes: accountID.UUID(),
			Valid: true,
		},
		DeletedAt: pgtype.Timestamptz{
			Time: now,
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamptz{
			Time: now,
			Valid: true,
		},
	})
}