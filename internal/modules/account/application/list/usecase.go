package list

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/account/ports"
	"github.com/zchelalo/expense-control-back/internal/shared/idgen"
	"go.uber.org/zap"
)

type UseCase struct {
	accounts ports.AccountRepository
	users ports.UserRepository
	ids      idgen.Generator
	pagLimitDefault int
}

func New(
	accounts ports.AccountRepository,
	users ports.UserRepository,
	ids idgen.Generator,
	pagLimitDefault int,
) *UseCase {
	return &UseCase{
		accounts: accounts,
		users: users,
		ids:      ids,
		pagLimitDefault: pagLimitDefault,
	}
}

type Result struct {
	Account []domain.Account
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)

	if cmd.UserID.String() == "" {
		log.Warn("missing user ID in list accounts request",
			zap.String("stage", "validate_input"),
		)
		return Result{}, domain.ErrInvalidUserID
	}

	if cmd.Limit <= 0 {
		log.Warn("invalid pagination limit in list accounts request",
			zap.String("stage", "validate_input"),
			zap.Int("limit", cmd.Limit),
		)
		cmd.Limit = uc.pagLimitDefault
	}

	if cmd.CreatedAt == nil && cmd.AccountID != nil {
		log.Warn("account ID provided without created at in list accounts request",
			zap.String("stage", "validate_input"),
			zap.String("account_id", cmd.AccountID.String()),
		)
		return Result{}, ErrAccountIDWithoutCreatedAt
	}

	if cmd.CreatedAt != nil && cmd.AccountID == nil {
		log.Warn("created at provided without account ID in list accounts request",
			zap.String("stage", "validate_input"),
			zap.Time("created_at", *cmd.CreatedAt),
		)
		return Result{}, ErrCreatedAtWithoutAccountID
	}

	// Verify that the user exists
	exists, err := uc.users.Exists(ctx, cmd.UserID)
	if err != nil {
		log.Error("failed to check if user exists",
			zap.String("stage", "check_user_exists"),
			zap.Error(err),
		)
		return Result{}, err
	}
	if !exists {
		log.Warn("user not found for create account request",
			zap.String("stage", "check_user_exists"),
			zap.String("user_id", cmd.UserID.String()),
		)
		return Result{}, ports.ErrNotFound{Name: "user"}
	}

	accounts, err := uc.accounts.ListByUserID(ctx, cmd.UserID, cmd.CreatedAt, cmd.AccountID, cmd.Limit, cmd.IsBefore)
	if err != nil {
		log.Error("failed to list accounts",
			zap.String("stage", "list_accounts"),
			zap.Error(err),
		)
		return Result{}, err
	}

	return Result{Account: accounts}, nil
}