package list

import "errors"

var ErrCreatedAtWithoutMovementID = errors.New("created_at cannot be used without movement_id")
var ErrMovementIDWithoutCreatedAt = errors.New("movement_id cannot be used without created_at")
var ErrDateFromAfterDateTo = errors.New("date_from cannot be later than date_to")
