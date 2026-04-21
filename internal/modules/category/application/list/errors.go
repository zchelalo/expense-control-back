package list

import "errors"

var ErrCreatedAtWithoutCategoryID = errors.New("created_at cannot be used without category_id")
var ErrCategoryIDWithoutCreatedAt = errors.New("category_id cannot be used without created_at")
