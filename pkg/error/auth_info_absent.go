package error

import "fmt"

type AuthInfoAbsentError struct {
	msg string
	Err error
}

var AuthInfoAbsentErr *AuthInfoAbsentError

func NewAuthInfoAbsentError(msg string, err error) *AuthInfoAbsentError {
	return &AuthInfoAbsentError{
		msg: msg,
		Err: err,
	}
}

// Error

func (a *AuthInfoAbsentError) Error() string {
	switch {
	case a.msg != "" && a.Err != nil:
		return fmt.Sprintf("auth user info absent with msg [%s] and err [%v]", a.msg, a.Err)
	case a.msg == "":
		return fmt.Sprintf("auth user info absent with msg [%s]", a.msg)
	case a.Err != nil:
		return fmt.Sprintf("auth user info absent with err [%v]", a.Err)
	}

	return "auth user info absent"
}

// ========

func (a *AuthInfoAbsentError) Unwrap() error {
	return a.Err
}
