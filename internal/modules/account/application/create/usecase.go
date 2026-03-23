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

	if cmd.UserID.String() == "" {
		log.Warn("missing user ID in create account request",
			zap.String("stage", "validate_input"),
		)
		return Result{}, domain.ErrInvalidUserID
	}

	if cmd.Name.String() == "" {
		log.Warn("missing account name in create account request",
			zap.String("stage", "validate_input"),
		)
		return Result{}, domain.ErrInvalidName
	}

	if cmd.Balance.Float64() < 0 {
		log.Warn("invalid account balance in create account request",
			zap.String("stage", "validate_input"),
			zap.Float64("balance", cmd.Balance.Float64()),
		)
		return Result{}, domain.ErrInvalidBalance
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
	userID, err := domain.NewUserID(cmd.UserID.UUID())
	if err != nil {
		log.Error("failed to generate user ID",
			zap.String("stage", "generate_user_id"),
			zap.Error(err),
		)
		return Result{}, err
	}
	name, err := domain.NewName(cmd.Name.String())
	if err != nil {
		log.Error("failed to generate account name",
			zap.String("stage", "generate_account_name"),
			zap.Error(err),
		)
		return Result{}, err
	}
	balance, err := domain.NewBalance(cmd.Balance.Float64())
	if err != nil {
		log.Error("failed to generate account balance",
			zap.String("stage", "generate_account_balance"),
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