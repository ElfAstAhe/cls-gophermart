package error

import (
	"fmt"
)

type ModelValidationError struct {
	Model   string
	Message string
	Err     error
}

var ModelValidationErr *ModelValidationError

func NewModelValidationError(model string, message string, err error) *ModelValidationError {
	return &ModelValidationError{
		Model:   model,
		Message: message,
		Err:     err,
	}
}

func (m *ModelValidationError) Error() string {
	if m.Err != nil {
		return fmt.Sprintf("model [%s] validateion error with message [%s], err [%v]", m.Model, m.Message, m.Err)
	}

	return fmt.Sprintf("model [%s] validateion error with message [%s]", m.Model, m.Message)
}

func (m *ModelValidationError) Unwrap() error {
	return m.Err
}
