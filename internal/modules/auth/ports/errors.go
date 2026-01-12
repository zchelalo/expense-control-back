package ports

import "fmt"

type ErrAlreadyExists struct {
	Name string
}

func (e ErrAlreadyExists) Error() string {
	return fmt.Sprintf("%s already exists", e.Name)
}

type ErrNotFound struct {
	Name string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found", e.Name)
}