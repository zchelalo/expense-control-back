package list

import "errors"

var ErrCreatedAtWithoutAccountID = errors.New("created_at cannot be used without account_id")
var ErrAccountIDWithoutCreatedAt = errors.New("account_id cannot be used without created_at")