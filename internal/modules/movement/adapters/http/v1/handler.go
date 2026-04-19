package v1

import (
	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/byid"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/create"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/delete"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/application/list"
)

type Handler struct {
	createUC *create.UseCase
	listUC   *list.UseCase
	byIDUC   *byid.UseCase
	deleteUC *delete.UseCase
}

func NewHandler(
	createUC *create.UseCase,
	listUC *list.UseCase,
	byIDUC *byid.UseCase,
	deleteUC *delete.UseCase,
) *Handler {
	return &Handler{
		createUC: createUC,
		listUC:   listUC,
		byIDUC:   byIDUC,
		deleteUC: deleteUC,
	}
}
