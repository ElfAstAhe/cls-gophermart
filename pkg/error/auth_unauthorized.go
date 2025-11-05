package error

import (
	"fmt"
	"strings"
)

type AuthUnauthorizedError struct {
	message string
	err     error
}

var AuthUnauthorizedErr *AuthUnauthorizedError

func NewAuthUnauthorizedError(message string, err error) *AuthUnauthorizedError {
	return &AuthUnauthorizedError{
		message: message,
		err:     err,
	}
}

func (e *AuthUnauthorizedError) Error() string {
	if strings.TrimSpace(e.message) == "" {
		return "unauthorized"
	}
	if e.err != nil {
		return fmt.Sprintf("unauthorized: [%s] with error [%v]", e.message, e.err)
	}

	return fmt.Sprintf("unauthorized: [%s]", e.message)
}

func (e *AuthUnauthorizedError) Unwrap() error {
	return e.err
}
