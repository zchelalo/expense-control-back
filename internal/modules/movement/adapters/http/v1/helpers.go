package v1

import (
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
	"github.com/zchelalo/expense-control-back/internal/shared/localization"
)

func mapMovementDetails(details domain.MovementDetails, language string) movementResponse {
	movement := details.Movement()
	movementType := details.MovementType()
	category := details.Category()
	account := details.Account()
	categoryName := localizeCategoryName(category, language)
	description := movement.Description().String()
	if category.IsSystem() {
		description = categoryName
	}

	return movementResponse{
		ID:          movement.ID().String(),
		Amount:      movement.Amount().Float64(),
		Description: description,
		MovementType: movementTypeResponse{
			ID:   movementType.ID().String(),
			Key:  movementType.Key(),
			Name: movementType.Name(),
		},
		Category: categoryResponse{
			ID:        category.ID().String(),
			Name:      categoryName,
			IsSystem:  category.IsSystem(),
			SystemKey: category.SystemKey(),
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

func requestLanguage(acceptLanguage string) string {
	return localization.ResolveLanguage(acceptLanguage)
}

func localizeCategoryName(category domain.Category, language string) string {
	if !category.IsSystem() {
		return category.Name()
	}

	if localized, ok := localization.LocalizeSystemCategoryName(category.SystemKey(), language); ok {
		return localized
	}

	return category.Name()
}
