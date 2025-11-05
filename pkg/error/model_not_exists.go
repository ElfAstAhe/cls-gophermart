package error

import "fmt"

type ModelNotExistsError struct {
	Model string
	Key   any
}

var ModelNotExistsErr *ModelNotExistsError

func NewModelNotExistsError(model string, key any) *ModelNotExistsError {
	return &ModelNotExistsError{
		Model: model,
		Key:   key,
	}
}

func (e *ModelNotExistsError) Error() string {
	return fmt.Sprintf("model [%s] with key [%v] not exists", e.Model, e.Key)
}
