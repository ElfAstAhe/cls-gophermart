package utils

import (
	"strconv"

	errs "github.com/ElfAstAhe/cls-gophermart/pkg/error"
	"github.com/phedde/luhn-algorithm"
)

func ValidateOrderNumberByLuhn(orderNumber string) error {
	num, err := strconv.ParseInt(orderNumber, 10, 64)
	if err != nil {
		return errs.NewBllInvalidOrderNumberError(orderNumber)
	}
	if !(luhn.IsValid(num)) {
		return errs.NewBllInvalidOrderNumberError(orderNumber)
	}

	return nil
}
