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

	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		log.Warn("invalid user ID in list accounts request",
			zap.String("stage", "validate_input"),
			zap.String("user_id", cmd.UserID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	var accountID *domain.AccountID
	if cmd.AccountID != nil {
		id, err := domain.NewAccountID(*cmd.AccountID)
		if err != nil {
			log.Warn("invalid account ID in list accounts request",
				zap.String("stage", "validate_input"),
				zap.String("account_id", cmd.AccountID.String()),
				zap.Error(err),
			)
			return Result{}, err
		}
		accountID = &id
	}

	if cmd.Limit <= 0 {
		log.Warn("invalid pagination limit in list accounts request",
			zap.String("stage", "validate_input"),
			zap.Int("limit", cmd.Limit),
		)
		cmd.Limit = uc.pagLimitDefault
	}

	if cmd.CreatedAt == nil && accountID != nil {
		log.Warn("account ID provided without created at in list accounts request",
			zap.String("stage", "validate_input"),
			zap.String("account_id", accountID.String()),
		)
		return Result{}, ErrAccountIDWithoutCreatedAt
	}

	if cmd.CreatedAt != nil && accountID == nil {
		log.Warn("created at provided without account ID in list accounts request",
			zap.String("stage", "validate_input"),
			zap.Time("created_at", *cmd.CreatedAt),
		)
		return Result{}, ErrCreatedAtWithoutAccountID
	}

	// Verify that the user exists
	exists, err := uc.users.Exists(ctx, userID)
	if err != nil {
		log.Error("failed to check if user exists",
			zap.String("stage", "check_user_exists"),
			zap.Error(err),
		)
		return Result{}, err
	}
	if !exists {
		log.Warn("user not found for list accounts request",
			zap.String("stage", "check_user_exists"),
			zap.String("user_id", userID.String()),
		)
		return Result{}, ports.ErrNotFound{Name: "user"}
	}

	accounts, err := uc.accounts.ListByUserID(ctx, userID, cmd.Name, cmd.CreatedAt, accountID, cmd.Limit, cmd.IsBefore)
	if err != nil {
		log.Error("failed to list accounts",
			zap.String("stage", "list_accounts"),
			zap.Error(err),
		)
		return Result{}, err
	}

	return Result{Account: accounts}, nil
}