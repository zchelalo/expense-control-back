package v1

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	uuidparse "github.com/zchelalo/expense-control-back/pkg/parse"
	"github.com/zchelalo/expense-control-back/pkg/response"
)

const simpleDateLayout = "2006-01-02"

type movementQueryFilters struct {
	AccountID      *uuid.UUID
	CategoryID     *uuid.UUID
	MovementTypeID *uuid.UUID
	DateFrom       *time.Time
	DateTo         *time.Time
}

func parseMovementQueryFilters(queries url.Values) (movementQueryFilters, *response.APIError) {
	accountID, err := uuidparse.OptionalUUID(queries.Get("account_id"))
	if err != nil {
		return movementQueryFilters{}, &response.APIError{
			Code:    "invalid_account_id",
			Message: "invalid account ID format",
		}
	}

	categoryID, err := uuidparse.OptionalUUID(queries.Get("category_id"))
	if err != nil {
		return movementQueryFilters{}, &response.APIError{
			Code:    "invalid_category_id",
			Message: "invalid category ID format",
		}
	}

	movementTypeID, err := uuidparse.OptionalUUID(queries.Get("movement_type_id"))
	if err != nil {
		return movementQueryFilters{}, &response.APIError{
			Code:    "invalid_movement_type_id",
			Message: "invalid movement type ID format",
		}
	}

	dateFrom, err := parseDateQuery(firstNonEmpty(queries.Get("date_from"), queries.Get("from_date")), false)
	if err != nil {
		return movementQueryFilters{}, &response.APIError{
			Code:    "invalid_date_from",
			Message: "date_from must be RFC3339 or YYYY-MM-DD",
		}
	}

	dateTo, err := parseDateQuery(firstNonEmpty(queries.Get("date_to"), queries.Get("to_date")), true)
	if err != nil {
		return movementQueryFilters{}, &response.APIError{
			Code:    "invalid_date_to",
			Message: "date_to must be RFC3339 or YYYY-MM-DD",
		}
	}

	if dateFrom != nil && dateTo != nil && dateFrom.After(*dateTo) {
		return movementQueryFilters{}, &response.APIError{
			Code:    "invalid_date_range",
			Message: "date_from cannot be later than date_to",
		}
	}

	return movementQueryFilters{
		AccountID:      accountID,
		CategoryID:     categoryID,
		MovementTypeID: movementTypeID,
		DateFrom:       dateFrom,
		DateTo:         dateTo,
	}, nil
}

func parseDateQuery(raw string, endOfDay bool) (*time.Time, error) {
	if raw == "" {
		return nil, nil
	}

	for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return &parsed, nil
		}
	}

	parsed, err := time.ParseInLocation(simpleDateLayout, raw, time.UTC)
	if err != nil {
		return nil, err
	}

	if endOfDay {
		parsed = parsed.Add(24*time.Hour - time.Nanosecond)
	}

	return &parsed, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}
