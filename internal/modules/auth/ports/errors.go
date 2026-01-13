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

type ErrTokenInvalid struct{
	Name string
}

func (e ErrTokenInvalid) Error() string {
	return fmt.Sprintf("%s token is invalid", e.Name)
}

type ErrTokenExpired struct{
	Name string
}

func (e ErrTokenExpired) Error() string {
	return fmt.Sprintf("%s token has expired", e.Name)
}