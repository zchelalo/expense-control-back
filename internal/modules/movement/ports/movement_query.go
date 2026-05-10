package ports

import (
	"context"
	"time"

	"github.com/zchelalo/expense-control-back/internal/modules/movement/domain"
)

type ListMovementsFilter struct {
	AccountID      *domain.AccountID
	CategoryID     *domain.CategoryID
	MovementTypeID *domain.MovementTypeID
	DateFrom       *time.Time
	DateTo         *time.Time
	CreatedAt      *time.Time
	MovementID     *domain.MovementID
	Limit          int
	IsBefore       bool
}

type StatsFilter struct {
	AccountID      *domain.AccountID
	CategoryID     *domain.CategoryID
	MovementTypeID *domain.MovementTypeID
	DateFrom       *time.Time
	DateTo         *time.Time
}

type MovementStatsOverview struct {
	TotalMovements int64
	IncomeCount    int64
	ExpenseCount   int64
	IncomeTotal    float64
	ExpenseTotal   float64
	NetTotal       float64
}

type MovementStatsByAccount struct {
	AccountID     domain.AccountID
	AccountName   string
	MovementCount int64
	IncomeCount   int64
	ExpenseCount  int64
	IncomeTotal   float64
	ExpenseTotal  float64
	NetTotal      float64
}

type MovementStatsByCategory struct {
	CategoryID        domain.CategoryID
	CategoryName      string
	CategoryIsSystem  bool
	CategorySystemKey string
	MovementCount     int64
	IncomeCount       int64
	ExpenseCount      int64
	IncomeTotal       float64
	ExpenseTotal      float64
	NetTotal          float64
}

// QueryRepository is the read side to hydrate details.
type QueryRepository interface {
	ByIDForUser(ctx context.Context, movementID domain.MovementID, userID domain.UserID) (domain.MovementDetails, error)
	ListByUserID(ctx context.Context, userID domain.UserID, filter ListMovementsFilter) ([]domain.MovementDetails, error)
	GetStatsOverviewByUserID(ctx context.Context, userID domain.UserID, filter StatsFilter) (MovementStatsOverview, error)
	ListStatsByAccountByUserID(ctx context.Context, userID domain.UserID, filter StatsFilter) ([]MovementStatsByAccount, error)
	ListStatsByCategoryByUserID(ctx context.Context, userID domain.UserID, filter StatsFilter) ([]MovementStatsByCategory, error)
}
