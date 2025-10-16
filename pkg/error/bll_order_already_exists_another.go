package error

import "fmt"

type BllOrderAlreadyExistsAnotherError struct {
	number string
}

var BllOrderAlreadyExistsAnotherErr *BllOrderAlreadyExistsAnotherError

func NewBllOrderAlreadyExistsAnotherError(number string) *BllOrderAlreadyExistsAnotherError {
	return &BllOrderAlreadyExistsAnotherError{
		number: number,
	}
}

func (e *BllOrderAlreadyExistsAnotherError) Error() string {
	return fmt.Sprintf("order num [%s] by another user already exists", e.number)
}
