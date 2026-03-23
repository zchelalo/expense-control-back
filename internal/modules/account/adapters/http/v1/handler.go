package v1

import (
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/create"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/list"
)

type Handler struct {
	createUC *create.UseCase
	listUC   *list.UseCase
}

func NewHandler(
	createUC *create.UseCase,
	listUC *list.UseCase,
) *Handler {
	return &Handler{
		createUC: createUC,
		listUC:   listUC,
	}
}