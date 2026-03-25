package updatename

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/account/ports"
	"github.com/zchelalo/expense-control-back/internal/shared/clock"
	"go.uber.org/zap"
)

type UseCase struct {
	accounts ports.AccountRepository
	users ports.UserRepository
	clock    clock.Clock
}

func New(
	accounts ports.AccountRepository,
	users ports.UserRepository,
	clock clock.Clock,
) *UseCase {
	return &UseCase{
		accounts: accounts,
		users: users,
		clock: clock,
	}
}

type Result struct {
	Account domain.Account
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)

	if cmd.UserID.String() == "" {
		log.Warn("missing user ID in get by id account request",
			zap.String("stage", "validate_input"),
		)
		return Result{}, domain.ErrInvalidUserID
	}

	if cmd.AccountID.String() == "" {
		log.Warn("missing account ID in get by id account request",
			zap.String("stage", "validate_input"),
		)
		return Result{}, domain.ErrInvalidAccountID
	}

	if cmd.Name.String() == "" {
		log.Warn("missing name in update account name request",
			zap.String("stage", "validate_input"),
		)
		return Result{}, domain.ErrInvalidName
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

	account, err := uc.accounts.ByID(ctx, cmd.AccountID)
	if err != nil {
		log.Error("failed to get account by ID",
			zap.String("stage", "get_account_by_id"),
			zap.Error(err),
		)
		return Result{}, err
	}

	if account.UserID().String() != cmd.UserID.String() {
		log.Warn("account doesn't belong to user in get by id account request",
			zap.String("stage", "validate_account_belongs_to_user"),
			zap.String("account_id", cmd.AccountID.String()),
			zap.String("user_id", cmd.UserID.String()),
		)
		return Result{}, ErrAccountDoesntBelongToUser
	}

	now := uc.clock.Now()
	err = uc.accounts.UpdateName(ctx, cmd.AccountID, cmd.Name, now)
	if err != nil {
		log.Error("failed to update account name",
			zap.String("stage", "update_account_name"),
			zap.Error(err),
		)
		return Result{}, err
	}
	
	return Result{Account: domain.RehydrateAccount(
		account.ID(),
		cmd.Name,
		account.Balance(),
		account.UserID(),
		account.CreatedAt(),
		now,
		account.DeletedAt(),
	)}, nil
}