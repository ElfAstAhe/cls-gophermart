package error

import "fmt"

type AuthPasswordIncorrectError struct {
	username string
}

var AuthPasswordIncorrectErr *AuthPasswordIncorrectError

func NewAuthPasswordIncorrectError(username string) *AuthPasswordIncorrectError {
	return &AuthPasswordIncorrectError{
		username: username,
	}
}

func (aic *AuthPasswordIncorrectError) Error() string {
	return fmt.Sprintf("password incorrect for user [%s]", aic.username)
}
