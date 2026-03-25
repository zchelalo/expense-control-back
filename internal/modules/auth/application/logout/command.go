package logout

import "github.com/google/uuid"

type Command struct {
	SubjectID    uuid.UUID
	RefreshToken string
}
