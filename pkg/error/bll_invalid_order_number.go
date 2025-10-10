package error

import "fmt"

type BllInvalidOrderNumberError struct {
	Number string
}

var BllInvalidOrderNumberErr *BllInvalidOrderNumberError

func NewBllInvalidOrderNumberError(number string) *BllInvalidOrderNumberError {
	return &BllInvalidOrderNumberError{
		Number: number,
	}
}

func (e *BllInvalidOrderNumberError) Error() string {
	return fmt.Sprintf("order number [%s] is invalid", e.Number)
}
