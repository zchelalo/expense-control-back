package ports

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/modules/movementtype/domain"
)

type MovementTypeRepository interface {
	List(ctx context.Context) ([]domain.MovementType, error)
}
