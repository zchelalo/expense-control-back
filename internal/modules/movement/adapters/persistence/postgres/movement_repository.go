package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	movementdb "github.com/zchelalo/expense-control-back/internal/db/sqlc/movement"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/ports"
	pgutil "github.com/zchelalo/expense-control-back/internal/shared/postgresutil"
)

type movementTxDB interface {
	movementdb.DBTX
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type MovementRepo struct {
	db movementTxDB
	q  *movementdb.Queries
}

func NewMovementRepo(db movementTxDB) *MovementRepo {
	return &MovementRepo{
		db: db,
		q:  movementdb.New(db),
	}
}

func (r *MovementRepo) Create(ctx context.Context, movement domain.Movement, operation ports.BalanceOperation) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)

	if err := r.applyBalanceOperation(ctx, q, movement, operation, movement.UpdatedAt()); err != nil {
		return err
	}

	amount, err := pgutil.NumericFromFloat64(movement.Amount().Float64())
	if err != nil {
		return err
	}

	if err := q.CreateMovement(ctx, movementdb.CreateMovementParams{
		ID:             pgutil.UUID(movement.ID()),
		Amount:         amount,
		Description:    movement.Description().String(),
		MovementTypeID: pgutil.UUID(movement.MovementTypeID()),
		CategoryID:     pgutil.UUID(movement.CategoryID()),
		AccountID:      pgutil.UUID(movement.AccountID()),
		UserID:         pgutil.UUID(movement.UserID()),
		CreatedAt:      pgutil.Timestamptz(movement.CreatedAt()),
		UpdatedAt:      pgutil.Timestamptz(movement.UpdatedAt()),
		DeletedAt:      pgutil.OptionalTimestamptz(movement.DeletedAt()),
	}); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *MovementRepo) ByID(ctx context.Context, movementID domain.MovementID) (domain.Movement, error) {
	movement, err := r.q.GetMovementByID(ctx, pgutil.UUID(movementID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Movement{}, ports.ErrNotFound{Name: "movement"}
		}
		return domain.Movement{}, err
	}

	return hydrateMovement(
		movement.ID,
		movement.Amount,
		movement.Description,
		movement.MovementTypeID,
		movement.CategoryID,
		movement.AccountID,
		movement.UserID,
		movement.CreatedAt,
		movement.UpdatedAt,
		movement.DeletedAt,
	)
}

func (r *MovementRepo) Delete(ctx context.Context, movement domain.Movement, operation ports.BalanceOperation, now time.Time) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)

	if err := r.applyBalanceOperation(ctx, q, movement, operation, now); err != nil {
		return err
	}

	affectedRows, err := q.DeleteMovement(ctx, movementdb.DeleteMovementParams{
		ID:        pgutil.UUID(movement.ID()),
		DeletedAt: pgutil.Timestamptz(now),
		UpdatedAt: pgutil.Timestamptz(now),
	})
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return ports.ErrNotFound{Name: "movement"}
	}

	return tx.Commit(ctx)
}

func (r *MovementRepo) applyBalanceOperation(ctx context.Context, q *movementdb.Queries, movement domain.Movement, operation ports.BalanceOperation, now time.Time) error {
	amount, err := pgutil.NumericFromFloat64(movement.Amount().Float64())
	if err != nil {
		return err
	}

	var affectedRows int64
	switch operation {
	case ports.BalanceOperationCredit:
		affectedRows, err = q.IncreaseAccountBalance(ctx, movementdb.IncreaseAccountBalanceParams{
			Balance:   amount,
			ID:        pgutil.UUID(movement.AccountID()),
			UserID:    pgutil.UUID(movement.UserID()),
			UpdatedAt: pgutil.Timestamptz(now),
		})
	case ports.BalanceOperationDebit:
		affectedRows, err = q.DecreaseAccountBalance(ctx, movementdb.DecreaseAccountBalanceParams{
			Balance:   amount,
			ID:        pgutil.UUID(movement.AccountID()),
			UserID:    pgutil.UUID(movement.UserID()),
			UpdatedAt: pgutil.Timestamptz(now),
		})
	default:
		return domain.ErrInvalidMovementType
	}
	if err != nil {
		return err
	}
	if affectedRows > 0 {
		return nil
	}

	accountExists, err := q.AccountExistsByUserID(ctx, movementdb.AccountExistsByUserIDParams{
		ID:     pgutil.UUID(movement.AccountID()),
		UserID: pgutil.UUID(movement.UserID()),
	})
	if err != nil {
		return err
	}
	if !accountExists {
		return ports.ErrNotFound{Name: "account"}
	}

	return domain.ErrInsufficientAccountBalance
}
