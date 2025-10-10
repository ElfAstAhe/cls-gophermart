package utils

import (
	"strconv"

	_err "github.com/ElfAstAhe/cls-gophermart/pkg/error"
	"github.com/phedde/luhn-algorithm"
)

func ValidateOrderNumberByLuhn(orderNumber string) error {
	num, err := strconv.ParseInt(orderNumber, 10, 64)
	if err != nil {
		return _err.NewBllInvalidOrderNumberError(orderNumber)
	}
	if !(luhn.IsValid(num)) {
		return _err.NewBllInvalidOrderNumberError(orderNumber)
	}

	return nil
}
