package create

import (
	"context"

	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/account/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/account/ports"
	"github.com/zchelalo/expense-control-back/internal/shared/clock"
	"github.com/zchelalo/expense-control-back/internal/shared/idgen"
	"go.uber.org/zap"
)

type UseCase struct {
	accounts ports.AccountRepository
	users ports.UserRepository
	clock    clock.Clock
	ids      idgen.Generator
}

func New(
	accounts ports.AccountRepository,
	users ports.UserRepository,
	clock clock.Clock,
	ids idgen.Generator,
) *UseCase {
	return &UseCase{
		accounts: accounts,
		users: users,
		clock:    clock,
		ids:      ids,
	}
}

type Result struct {
	Account *domain.Account
}

func (uc *UseCase) Execute(ctx context.Context, cmd Command) (Result, error) {
	log := middleware.LoggerFrom(ctx)

	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		log.Warn("invalid user ID in create account request",
			zap.String("stage", "validate_input"),
			zap.String("user_id", cmd.UserID.String()),
			zap.Error(err),
		)
		return Result{}, err
	}

	name, err := domain.NewName(cmd.Name)
	if err != nil {
		log.Warn("invalid account name in create account request",
			zap.String("stage", "validate_input"),
			zap.String("name", cmd.Name),
			zap.Error(err),
		)
		return Result{}, err
	}

	balance, err := domain.NewBalance(cmd.Balance)
	if err != nil {
		log.Warn("invalid account balance in create account request",
			zap.String("stage", "validate_input"),
			zap.Float64("balance", cmd.Balance),
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
		log.Warn("user not found for create account request",
			zap.String("stage", "check_user_exists"),
			zap.String("user_id", userID.String()),
		)
		return Result{}, ports.ErrNotFound{Name: "user"}
	}

	// Create the account
	now := uc.clock.Now()

	accountID, err := domain.NewAccountID(uc.ids.NewUUID())
	if err != nil {
		log.Error("failed to generate account ID",
			zap.String("stage", "generate_account_id"),
			zap.Error(err),
		)
		return Result{}, err
	}

	account := domain.NewAccount(
		accountID,
		name,
		balance,
		userID,
		now,
	)
	if err := uc.accounts.Create(ctx, account); err != nil {
		log.Error("failed to create account",
			zap.String("stage", "create_account"),
			zap.Error(err),
		)
		return Result{}, err
	}

	return Result{Account: &account}, nil
}