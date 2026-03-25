package delete

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
	clock clock.Clock
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
	Success bool
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)

	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		log.Warn("invalid user ID in delete account request",
			zap.String("stage", "validate_input"),
			zap.String("user_id", cmd.UserID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	accountID, err := domain.NewAccountID(cmd.AccountID)
	if err != nil {
		log.Warn("invalid account ID in delete account request",
			zap.String("stage", "validate_input"),
			zap.String("account_id", cmd.AccountID.String()),
			zap.Error(err),
		)
		return Result{}, err
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
		log.Warn("user not found for delete account request",
			zap.String("stage", "check_user_exists"),
			zap.String("user_id", userID.String()),
		)
		return Result{}, ports.ErrNotFound{Name: "user"}
	}

	account, err := uc.accounts.ByID(ctx, accountID)
	if err != nil {
		log.Error("failed to get account by ID",
			zap.String("stage", "get_account_by_id"),
			zap.Error(err),
		)
		return Result{}, err
	}

	if account.UserID().String() != userID.String() {
		log.Warn("account doesn't belong to user in delete account request",
			zap.String("stage", "validate_account_belongs_to_user"),
			zap.String("account_id", accountID.String()),
			zap.String("user_id", userID.String()),
		)
		return Result{}, ErrAccountDoesntBelongToUser
	}

	now := uc.clock.Now()
	err = uc.accounts.Delete(ctx, accountID, now)
	if err != nil {
		log.Error("failed to delete account",
			zap.String("stage", "delete_account"),
			zap.Error(err),
		)
		return Result{}, err
	}

	return Result{Success: true}, nil
}