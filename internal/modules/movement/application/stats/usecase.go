package stats

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zchelalo/expense-control-back/internal/middleware"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
	"github.com/zchelalo/expense-control-back/internal/modules/movement/ports"
	"go.uber.org/zap"
)

var ErrDateFromAfterDateTo = errors.New("date_from cannot be later than date_to")

type UseCase struct {
	query ports.QueryRepository
	users ports.UserReferenceRepository
}

func New(
	query ports.QueryRepository,
	users ports.UserReferenceRepository,
) *UseCase {
	return &UseCase{
		query: query,
		users: users,
	}
}

type OverviewResult struct {
	Overview ports.MovementStatsOverview
}

type ByAccountResult struct {
	Accounts []ports.MovementStatsByAccount
}

type ByCategoryResult struct {
	Categories []ports.MovementStatsByCategory
}

func (uc *UseCase) Overview(ctx context.Context, cmd Command) (OverviewResult, error) {
	userID, filter, err := uc.validate(ctx, cmd, "get movement stats overview request")
	if err != nil {
		return OverviewResult{}, err
	}

	overview, err := uc.query.GetStatsOverviewByUserID(ctx, userID, filter)
	if err != nil {
		middleware.LoggerFrom(ctx).Error("failed to get movement stats overview",
			zap.String("stage", "get_movement_stats_overview"),
			zap.Error(err),
		)
		return OverviewResult{}, err
	}

	return OverviewResult{Overview: overview}, nil
}

func (uc *UseCase) ByAccount(ctx context.Context, cmd Command) (ByAccountResult, error) {
	userID, filter, err := uc.validate(ctx, cmd, "list movement stats by account request")
	if err != nil {
		return ByAccountResult{}, err
	}

	accounts, err := uc.query.ListStatsByAccountByUserID(ctx, userID, filter)
	if err != nil {
		middleware.LoggerFrom(ctx).Error("failed to list movement stats by account",
			zap.String("stage", "list_movement_stats_by_account"),
			zap.Error(err),
		)
		return ByAccountResult{}, err
	}

	return ByAccountResult{Accounts: accounts}, nil
}

func (uc *UseCase) ByCategory(ctx context.Context, cmd Command) (ByCategoryResult, error) {
	userID, filter, err := uc.validate(ctx, cmd, "list movement stats by category request")
	if err != nil {
		return ByCategoryResult{}, err
	}

	categories, err := uc.query.ListStatsByCategoryByUserID(ctx, userID, filter)
	if err != nil {
		middleware.LoggerFrom(ctx).Error("failed to list movement stats by category",
			zap.String("stage", "list_movement_stats_by_category"),
			zap.Error(err),
		)
		return ByCategoryResult{}, err
	}

	return ByCategoryResult{Categories: categories}, nil
}

func (uc *UseCase) validate(ctx context.Context, cmd Command, action string) (domain.UserID, ports.StatsFilter, error) {
	log := middleware.LoggerFrom(ctx)

	userID, err := domain.NewUserID(cmd.UserID)
	if err != nil {
		log.Warn("invalid user ID in "+action,
			zap.String("stage", "validate_input"),
			zap.String("user_id", cmd.UserID.String()),
			zap.Error(err),
		)
		return domain.UserID{}, ports.StatsFilter{}, err
	}

	accountID, err := optionalAccountID(cmd.AccountID)
	if err != nil {
		log.Warn("invalid account ID in "+action,
			zap.String("stage", "validate_input"),
			zap.String("account_id", cmd.AccountID.String()),
			zap.Error(err),
		)
		return domain.UserID{}, ports.StatsFilter{}, err
	}

	categoryID, err := optionalCategoryID(cmd.CategoryID)
	if err != nil {
		log.Warn("invalid category ID in "+action,
			zap.String("stage", "validate_input"),
			zap.String("category_id", cmd.CategoryID.String()),
			zap.Error(err),
		)
		return domain.UserID{}, ports.StatsFilter{}, err
	}

	movementTypeID, err := optionalMovementTypeID(cmd.MovementTypeID)
	if err != nil {
		log.Warn("invalid movement type ID in "+action,
			zap.String("stage", "validate_input"),
			zap.String("movement_type_id", cmd.MovementTypeID.String()),
			zap.Error(err),
		)
		return domain.UserID{}, ports.StatsFilter{}, err
	}

	if cmd.DateFrom != nil && cmd.DateTo != nil && cmd.DateFrom.After(*cmd.DateTo) {
		log.Warn("invalid date range in "+action,
			zap.String("stage", "validate_input"),
			zap.Time("date_from", *cmd.DateFrom),
			zap.Time("date_to", *cmd.DateTo),
		)
		return domain.UserID{}, ports.StatsFilter{}, ErrDateFromAfterDateTo
	}

	exists, err := uc.users.Exists(ctx, userID)
	if err != nil {
		log.Error("failed to check if user exists",
			zap.String("stage", "check_user_exists"),
			zap.Error(err),
		)
		return domain.UserID{}, ports.StatsFilter{}, err
	}
	if !exists {
		log.Warn("user not found in "+action,
			zap.String("stage", "check_user_exists"),
			zap.String("user_id", userID.String()),
		)
		return domain.UserID{}, ports.StatsFilter{}, ports.ErrNotFound{Name: "user"}
	}

	return userID, ports.StatsFilter{
		AccountID:      accountID,
		CategoryID:     categoryID,
		MovementTypeID: movementTypeID,
		DateFrom:       cmd.DateFrom,
		DateTo:         cmd.DateTo,
	}, nil
}

func optionalAccountID(raw *uuid.UUID) (*domain.AccountID, error) {
	if raw == nil {
		return nil, nil
	}

	id, err := domain.NewAccountID(*raw)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func optionalCategoryID(raw *uuid.UUID) (*domain.CategoryID, error) {
	if raw == nil {
		return nil, nil
	}

	id, err := domain.NewCategoryID(*raw)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func optionalMovementTypeID(raw *uuid.UUID) (*domain.MovementTypeID, error) {
	if raw == nil {
		return nil, nil
	}

	id, err := domain.NewMovementTypeID(*raw)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
