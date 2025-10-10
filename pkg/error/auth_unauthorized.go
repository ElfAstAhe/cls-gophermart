package error

import (
	"fmt"
	"strings"
)

type AuthUnauthorizedError struct {
	message string
}

var AuthUnauthorizedErr *AuthUnauthorizedError

func NewAuthUnauthorizedError(message string) *AuthUnauthorizedError {
	return &AuthUnauthorizedError{
		message: message,
	}
}

func (e *AuthUnauthorizedError) Error() string {
	if strings.TrimSpace(e.message) == "" {
		return "unauthorized"
	}

	return fmt.Sprintf("unauthorized with message [%s]", e.message)
}
