package postgres

import (
	"context"

	movementdb "github.com/zchelalo/expense-control-back/internal/db/sqlc/movement"
	"github.com/zchelalo/expense-control-back/internal/modules/movementtype/domain"
)

type MovementTypeRepo struct {
	q *movementdb.Queries
}

func NewMovementTypeRepo(db movementdb.DBTX) *MovementTypeRepo {
	return &MovementTypeRepo{q: movementdb.New(db)}
}

func (r *MovementTypeRepo) List(ctx context.Context) ([]domain.MovementType, error) {
	rows, err := r.q.ListMovementTypes(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.MovementType, len(rows))
	for i, row := range rows {
		id, err := domain.NewMovementTypeID(row.ID.Bytes)
		if err != nil {
			return nil, err
		}

		description := ""
		if row.Description.Valid {
			description = row.Description.String
		}

		movementType, err := domain.RehydrateMovementType(
			id,
			row.Key,
			row.Name,
			description,
		)
		if err != nil {
			return nil, err
		}

		result[i] = movementType
	}

	return result, nil
}
