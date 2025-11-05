package error

import "fmt"

type ModelAlreadyExistsError struct {
	Model string
	Key   any
}

var ModelAlreadyExistsErr *ModelAlreadyExistsError

func NewModelAlreadyExistsError(model string, key any) *ModelAlreadyExistsError {
	return &ModelAlreadyExistsError{
		Model: model,
		Key:   key,
	}
}

func (ae *ModelAlreadyExistsError) Error() string {
	return fmt.Sprintf("model [%s] with key [%v] already exists", ae.Model, ae.Key)
}
