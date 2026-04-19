package postgresutil

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UUIDValuer interface {
	UUID() uuid.UUID
}

func UUID(value UUIDValuer) pgtype.UUID {
	return pgtype.UUID{
		Bytes: value.UUID(),
		Valid: true,
	}
}

func OptionalUUID[T UUIDValuer](value *T) pgtype.UUID {
	if value == nil {
		return pgtype.UUID{Valid: false}
	}

	return UUID(*value)
}

func Timestamptz(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  value,
		Valid: true,
	}
}

func OptionalTimestamptz(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{Valid: false}
	}

	return Timestamptz(*value)
}

func TimestamptzPtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}

	t := value.Time
	return &t
}

func NumericFromFloat64(value float64) (pgtype.Numeric, error) {
	var numeric pgtype.Numeric
	if err := numeric.Scan(fmt.Sprintf("%f", value)); err != nil {
		return pgtype.Numeric{}, err
	}

	return numeric, nil
}

func NumericToFloat64(value pgtype.Numeric) (float64, error) {
	parsed, err := value.Float64Value()
	if err != nil {
		return 0, err
	}

	return parsed.Float64, nil
}
