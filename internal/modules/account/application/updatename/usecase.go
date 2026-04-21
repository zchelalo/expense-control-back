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
	users    ports.UserReferenceRepository
	clock    clock.Clock
}

func New(
	accounts ports.AccountRepository,
	users ports.UserReferenceRepository,
	clock clock.Clock,
) *UseCase {
	return &UseCase{
		accounts: accounts,
		users:    users,
		clock:    clock,
	}
}

type Result struct {
	Account domain.Account
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)

	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		log.Warn("invalid user ID in update account name request",
			zap.String("stage", "validate_input"),
			zap.String("user_id", cmd.UserID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	accountID, err := domain.NewAccountID(cmd.AccountID)
	if err != nil {
		log.Warn("invalid account ID in update account name request",
			zap.String("stage", "validate_input"),
			zap.String("account_id", cmd.AccountID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	name, err := domain.NewName(cmd.Name)
	if err != nil {
		log.Warn("invalid name in update account name request",
			zap.String("stage", "validate_input"),
			zap.String("name", cmd.Name),
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
		log.Warn("user not found for update account name request",
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
		log.Warn("account doesn't belong to user in update account name request",
			zap.String("stage", "validate_account_belongs_to_user"),
			zap.String("account_id", accountID.String()),
			zap.String("user_id", userID.String()),
		)
		return Result{}, ErrAccountDoesntBelongToUser
	}

	now := uc.clock.Now()
	err = uc.accounts.UpdateName(ctx, accountID, name, now)
	if err != nil {
		log.Error("failed to update account name",
			zap.String("stage", "update_account_name"),
			zap.Error(err),
		)
		return Result{}, err
	}

	return Result{Account: domain.RehydrateAccount(
		account.ID(),
		name,
		account.Balance(),
		account.UserID(),
		account.CreatedAt(),
		now,
		account.DeletedAt(),
	)}, nil
}
