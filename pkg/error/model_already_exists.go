package error

import "fmt"

type ModelAlreadyExistsError struct {
	Model string
}

var ModelAlreadyExistsErr *ModelAlreadyExistsError

func NewModelAlreadyExistsError(model string) *ModelAlreadyExistsError {
	return &ModelAlreadyExistsError{
		Model: model,
	}
}

func (err *ModelAlreadyExistsError) Error() string {
	return fmt.Sprintf("model [%s] already exists", err.Model)
}
