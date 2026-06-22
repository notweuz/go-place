package errs

import "fmt"

type AppError struct {
	Kind       error
	Err        error
	StatusCode int
	Message    string
}

func NewAppError(kind error, err error, statusCode int, message string) *AppError {
	return &AppError{
		Kind:       kind,
		Err:        err,
		StatusCode: statusCode,
		Message:    message,
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s (%d)", e.Kind, e.Message, e.StatusCode)
}

func (e *AppError) Unwrap() error {
	return e.Kind
}
