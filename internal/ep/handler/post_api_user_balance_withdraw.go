package handler

import (
	"errors"
	"net/http"

	_utl "github.com/ElfAstAhe/cls-gophermart/internal/utils"
	_err "github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

func (cr *AppChiRouter) postApiUserBalanceWithdraw(rw http.ResponseWriter, r *http.Request) {
	cr.log.Debug("postApiUserBalanceWithdraw start")
	defer cr.log.Debug("postApiUserBalanceWithdraw finish")

	defer _utl.CloseOnly(r.Body)

	if err := cr.usersFacade.CreateWithdraw(r.Context(), r.Body); err != nil {
		// 401
		if errors.As(err, &_err.AuthUnauthorizedErr) {
			http.Error(rw, err.Error(), http.StatusUnauthorized)

			return
		}
		// 402
		if errors.As(err, &_err.BllInsufficientBalanceErr) {
			http.Error(rw, err.Error(), http.StatusPaymentRequired)

			return
		}
		// 422
		if errors.As(err, &_err.BllInvalidOrderNumberErr) {
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
