package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	accountdb "github.com/zchelalo/expense-control-back/internal/db/sqlc/account"
	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/account/ports"
	pgutil "github.com/zchelalo/expense-control-back/internal/shared/postgresutil"
)

type AccountRepo struct {
	q *accountdb.Queries
}

func NewAccountRepo(db accountdb.DBTX) *AccountRepo {
	return &AccountRepo{q: accountdb.New(db)}
}

func (r *AccountRepo) Create(ctx context.Context, s domain.Account) error {
	balance, err := pgutil.NumericFromFloat64(s.Balance().Float64())
	if err != nil {
		return err
	}

	params := accountdb.CreateAccountParams{
		ID:        pgutil.UUID(s.ID()),
		Name:      s.Name().String(),
		Balance:   balance,
		UserID:    pgutil.UUID(s.UserID()),
		CreatedAt: pgutil.Timestamptz(s.CreatedAt()),
		UpdatedAt: pgutil.Timestamptz(s.UpdatedAt()),
		DeletedAt: pgutil.OptionalTimestamptz(s.DeletedAt()),
	}

	err = r.q.CreateAccount(ctx, params)
	return err
}

func (r *AccountRepo) ByID(ctx context.Context, id domain.AccountID) (domain.Account, error) {
	account, err := r.q.GetAccountByID(ctx, pgutil.UUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Account{}, ports.ErrNotFound{Name: "account"}
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

	balance, err := pgutil.NumericToFloat64(account.Balance)
	if err != nil {
		return domain.Account{}, err
	}

	parsedBalance, err := domain.NewBalance(balance)
	if err != nil {
		return domain.Account{}, err
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
		pgutil.TimestamptzPtr(account.DeletedAt),
	), nil
}
func (r *AccountRepo) ListByUserID(ctx context.Context, userID domain.UserID, name *string, createdAt *time.Time, accountID *domain.AccountID, limit int, isBefore bool) ([]domain.Account, error) {
	var accounts []accountdb.Account
	var err error

	var nameFilter pgtype.Text
	if name != nil {
		nameFilter = pgtype.Text{String: *name, Valid: true}
	} else {
		nameFilter = pgtype.Text{Valid: false}
	}

	if isBefore {
		accounts, err = r.q.ListAccountsByUserIDBefore(ctx, accountdb.ListAccountsByUserIDBeforeParams{
			UserID:  pgutil.UUID(userID),
			Column2: pgutil.Timestamptz(*createdAt),
			Column3: pgutil.UUID(*accountID),
			Limit:   int32(limit),
			Search:  nameFilter,
		})
	} else {
		params := accountdb.ListAccountsByUserIDAfterParams{
			UserID: pgutil.UUID(userID),
			Limit:  int32(limit),
			Search: nameFilter,
		}

		if createdAt != nil {
			params.Column2 = pgutil.Timestamptz(*createdAt)
		}

		if accountID != nil {
			params.Column3 = pgutil.UUID(*accountID)
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

		balance, err := pgutil.NumericToFloat64(account.Balance)
		if err != nil {
			return nil, err
		}

		parsedBalance, err := domain.NewBalance(balance)
		if err != nil {
			return nil, err
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
			pgutil.TimestamptzPtr(account.DeletedAt),
		)
	}

	return result, nil
}

func (r *AccountRepo) UpdateName(ctx context.Context, accountID domain.AccountID, name domain.Name, now time.Time) error {
	return r.q.UpdateAccountName(ctx, accountdb.UpdateAccountNameParams{
		ID:        pgutil.UUID(accountID),
		Name:      name.String(),
		UpdatedAt: pgutil.Timestamptz(now),
	})
}

func (r *AccountRepo) UpdateBalance(ctx context.Context, accountID domain.AccountID, balance domain.Balance, now time.Time) error {
	balanceNumeric, err := pgutil.NumericFromFloat64(balance.Float64())
	if err != nil {
		return err
	}

	return r.q.UpdateAccountBalance(ctx, accountdb.UpdateAccountBalanceParams{
		ID:        pgutil.UUID(accountID),
		Balance:   balanceNumeric,
		UpdatedAt: pgutil.Timestamptz(now),
	})
}

func (r *AccountRepo) Delete(ctx context.Context, accountID domain.AccountID, now time.Time) error {
	return r.q.DeleteAccount(ctx, accountdb.DeleteAccountParams{
		ID:        pgutil.UUID(accountID),
		DeletedAt: pgutil.Timestamptz(now),
		UpdatedAt: pgutil.Timestamptz(now),
	})
}
