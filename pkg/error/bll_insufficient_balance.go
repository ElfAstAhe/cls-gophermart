package error

import "fmt"

type BllInsufficientBalanceError struct {
	Current   float64
	Requested float64
}

var BllInsufficientBalanceErr *BllInsufficientBalanceError

func NewBllInsufficientBalanceError(current, requested float64) *BllInsufficientBalanceError {
	return &BllInsufficientBalanceError{
		Current:   current,
		Requested: requested,
	}
}

func (a *BllInsufficientBalanceError) Error() string {
	return fmt.Sprintf("insufficient balance, current [%.2f] requested [%.2f]", a.Current, a.Requested)
}
