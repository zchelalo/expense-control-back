package v1

import (
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/create"
)

type Handler struct {
	createUC 			*create.UseCase
}

func NewHandler(
	createUC *create.UseCase,
) *Handler {
	return &Handler{
		createUC: createUC,
	}
}