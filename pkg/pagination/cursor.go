package pagination

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func EncodeCursor(t time.Time, id uuid.UUID) string {
	s := fmt.Sprintf("%s,%s", t.Format(time.RFC3339Nano), id.String())
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func DecodeCursor(cursor string) (time.Time, uuid.UUID, error) {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}

	parts := strings.Split(string(b), ",")
	if len(parts) != 2 {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid cursor format")
	}

	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}

	id, err := uuid.Parse(parts[1])
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}

	return t, id, nil
}
