package list

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/movementtype/domain"
	"go.uber.org/zap"
)

type UseCase struct {
	movementTypes MovementTypeRepository
}

type MovementTypeRepository interface {
	List(ctx context.Context) ([]domain.MovementType, error)
}

func New(movementTypes MovementTypeRepository) *UseCase {
	return &UseCase{
		movementTypes: movementTypes,
	}
}

type Result struct {
	MovementTypes []domain.MovementType
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)
	movementTypes, err := uc.movementTypes.List(ctx)
	if err != nil {
		log.Error("failed to list movement types",
			zap.String("stage", "list_movement_types"),
			zap.Error(err),
		)
		return Result{}, err
	}

	return Result{MovementTypes: movementTypes}, nil
}
