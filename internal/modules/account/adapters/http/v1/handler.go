package v1

import (
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/byid"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/create"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/list"
	"github.com/zchelalo/expense-control-back/internal/modules/account/application/updatename"
)

type Handler struct {
	createUC *create.UseCase
	listUC   *list.UseCase
	byIDUC   *byid.UseCase
	updateNameUC *updatename.UseCase
}

func NewHandler(
	createUC *create.UseCase,
	listUC *list.UseCase,
	byIDUC *byid.UseCase,
	updateNameUC *updatename.UseCase,
) *Handler {
	return &Handler{
		createUC:   createUC,
		listUC:     listUC,
		byIDUC:     byIDUC,
		updateNameUC: updateNameUC,
	}
}