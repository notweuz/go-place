package errors

import "fmt"

type AppError struct {
	Kind    error
	Err     error
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Kind
}
