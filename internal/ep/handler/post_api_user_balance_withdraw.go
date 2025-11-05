package handler

import (
	"errors"
	"net/http"

	"github.com/ElfAstAhe/cls-gophermart/internal/utils"
	"github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

func (cr *AppChiRouter) postAPIUserBalanceWithdraw(rw http.ResponseWriter, r *http.Request) {
	cr.log.Debug("postAPIUserBalanceWithdraw start")
	defer cr.log.Debug("postAPIUserBalanceWithdraw finish")

	defer utils.CloseOnly(r.Body)

	if err := cr.usersFacade.CreateWithdraw(r.Context(), r.Body); err != nil {
		// 401
		if errors.As(err, &error.AuthUnauthorizedErr) {
			http.Error(rw, err.Error(), http.StatusUnauthorized)

			return
		}
		// 402
		if errors.As(err, &error.BllInsufficientBalanceErr) {
			http.Error(rw, err.Error(), http.StatusPaymentRequired)

			return
		}
		// 422
		if errors.As(err, &error.BllInvalidOrderNumberErr) {
			http.Error(rw, err.Error(), http.StatusUnprocessableEntity)

			return
		}

		// 500
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	// 200
	rw.WriteHeader(http.StatusOK)
}
