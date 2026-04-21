package v1

import "github.com/zchelalo/expense-control-back/internal/modules/movementtype/application/list"

type Handler struct {
	listUC *list.UseCase
}

func NewHandler(listUC *list.UseCase) *Handler {
	return &Handler{
		listUC: listUC,
	}
}
