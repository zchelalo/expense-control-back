package v1

import (
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/login"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/logout"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/refresh"
	"github.com/zchelalo/expense-control-back/internal/modules/auth/application/register"
)

type Handler struct {
	registerUC 		*register.UseCase
	loginUC 			*login.UseCase
	logoutUC 			*logout.UseCase
	refreshUC     *refresh.UseCase
	secureCookies bool
}

func NewHandler(
	registerUC *register.UseCase,
	loginUC *login.UseCase,
	logoutUC *logout.UseCase,
	refreshUC *refresh.UseCase,
	secureCookies bool,
) *Handler {
	return &Handler{
		registerUC: registerUC,
		loginUC: loginUC,
		logoutUC: logoutUC,
		refreshUC: refreshUC,
		secureCookies: secureCookies,
	}
}