package error

import "fmt"

type AuthCookieAbsentError struct {
	msg string
	Err error
}

var AuthCookieAbsent *AuthCookieAbsentError

func NewAuthCookieAbsentError(msg string, err error) *AuthCookieAbsentError {
	return &AuthCookieAbsentError{
		msg: msg,
		Err: err,
	}
}

func (a *AuthCookieAbsentError) Error() string {
	switch {
	case a.msg != "" && a.Err != nil:
		return fmt.Sprintf("auth cookie absent with msg [%s] and err [%v]", a.msg, a.Err)
	case a.msg == "":
		return fmt.Sprintf("auth cookie absent with msg [%s]", a.msg)
	case a.Err != nil:
		return fmt.Sprintf("auth cookie absent with err [%v]", a.Err)
	}

	return "auth cookie absent"
}

func (a *AuthCookieAbsentError) Unwrap() error {
	return a.Err
}
