package error

import "fmt"

type AppInvalidArgumentError struct {
	Param string
	Value any
}

var AppInvalidArgumentErr *AppInvalidArgumentError

func NewAppInvalidArgumentError(param string, value any) *AppInvalidArgumentError {
	return &AppInvalidArgumentError{Param: param, Value: value}
}

func (e *AppInvalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument [%s] with value [%v]", e.Param, e.Value)
}
