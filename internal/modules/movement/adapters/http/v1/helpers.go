package v1

import (
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
)

func mapMovementDetails(details domain.MovementDetails) movementResponse {
	movement := details.Movement()
	movementType := details.MovementType()
	category := details.Category()
	account := details.Account()

	return movementResponse{
		ID:          movement.ID().String(),
		Amount:      movement.Amount().Float64(),
		Description: movement.Description().String(),
		MovementType: movementTypeResponse{
			ID:   movementType.ID().String(),
			Key:  movementType.Key(),
			Name: movementType.Name(),
		},
		Category: categoryResponse{
			ID:   category.ID().String(),
			Name: category.Name(),
		},
		Account: accountResponse{
			ID:   account.ID().String(),
			Name: account.Name(),
		},
		UserID:    movement.UserID().String(),
		CreatedAt: movement.CreatedAt().UTC().Format(time.RFC3339),
		UpdatedAt: movement.UpdatedAt().UTC().Format(time.RFC3339),
	}
}
