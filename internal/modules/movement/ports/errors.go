package ports

import "fmt"

type ErrNotFound struct {
	Name string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found", e.Name)
}
