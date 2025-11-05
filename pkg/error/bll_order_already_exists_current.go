package error

import "fmt"

type BllOrderAlreadyExistsCurrentError struct {
	Number string
}

var BllOrderAlreadyExistsCurrentErr *BllOrderAlreadyExistsCurrentError

func NewBllOrderAlreadyExistsCurrentError(number string) *BllOrderAlreadyExistsCurrentError {
	return &BllOrderAlreadyExistsCurrentError{
		Number: number,
	}
}

func (e *BllOrderAlreadyExistsCurrentError) Error() string {
	return fmt.Sprintf("order num [%s] by current user already exists", e.Number)
}
