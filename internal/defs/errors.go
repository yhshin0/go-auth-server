package defs

import "fmt"

type CustomError struct {
	StatusCode int
	Code       string
	Message    string
	Cause      error
}

func (e CustomError) Error() string {
	if e.Cause == nil {
		return fmt.Sprintf("[%s] %s", e.Code, e.Message)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Cause.Error())
}

func (e CustomError) Wrap(err error) CustomError {
	e.Cause = err
	return e
}

func (e CustomError) Unwrap() error {
	return e.Cause
}
