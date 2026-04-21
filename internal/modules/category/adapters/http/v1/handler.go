package v1

import (
	"github.com/zchelalo/expense-control-back/internal/modules/category/application/create"
	"github.com/zchelalo/expense-control-back/internal/modules/category/application/delete"
	"github.com/zchelalo/expense-control-back/internal/modules/category/application/list"
)

type Handler struct {
	createUC *create.UseCase
	deleteUC *delete.UseCase
	listUC   *list.UseCase
}

func NewHandler(createUC *create.UseCase, deleteUC *delete.UseCase, listUC *list.UseCase) *Handler {
	return &Handler{
		createUC: createUC,
		deleteUC: deleteUC,
		listUC:   listUC,
	}
}
