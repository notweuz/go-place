package errs

import "fmt"

type AppError struct {
	Err        error
	StatusCode int
	Message    string
}

func NewAppError(err error, statusCode int, message string) *AppError {
	return &AppError{
		Err:        err,
		StatusCode: statusCode,
		Message:    message,
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s (%d)", e.Err, e.Message, e.StatusCode)
}

func (e *AppError) Unwrap() error {
	return e.Err
}
