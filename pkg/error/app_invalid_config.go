package error

import "fmt"

type AppInvalidConfigError struct {
	Name    string
	Value   any
	Message string
	Err     error
}

var AppInvalidConfigErr *AppInvalidConfigError

func NewAppInvalidConfigError(name string, value any, message string, err error) *AppInvalidConfigError {
	return &AppInvalidConfigError{
		Name:    name,
		Value:   value,
		Message: message,
		Err:     err,
	}
}

func NewAppInvalidConfigItemMsgError(name string, value any, message string) *AppInvalidConfigError {
	return NewAppInvalidConfigError(name, value, message, nil)
}

func NewAppInvalidConfigItemErr(name string, value any) *AppInvalidConfigError {
	return NewAppInvalidConfigItemMsgError(name, value, "")
}

func (apc *AppInvalidConfigError) Error() string {
	return fmt.Sprintf("invalid config item with name [%s] with value [%v] with message [%s]", apc.Name, apc.Value, apc.Message)
}

func (apc *AppInvalidConfigError) Unwrap() error {
	return apc.Err
}
