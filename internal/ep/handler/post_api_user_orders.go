package handler

import (
	"errors"
	"net/http"

	"github.com/ElfAstAhe/cls-gophermart/internal/utils"
	"github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

func (cr *AppChiRouter) postAPIUserOrders(rw http.ResponseWriter, r *http.Request) {
	cr.log.Debug("postAPIUserOrders start")
	defer cr.log.Debug("postAPIUserOrders finish")

	defer utils.CloseOnly(r.Body)

	// ToDo: make a table ? think
	if err := cr.usersFacade.CreateOrder(r.Context(), r.Body); err != nil {
		// 200
		if errors.As(err, &error.BllOrderAlreadyExistsCurrentErr) {
			rw.WriteHeader(http.StatusOK)

			return
		}
		// 400
		if errors.As(err, &error.ModelValidationErr) {
			http.Error(rw, err.Error(), http.StatusBadRequest)

			return
		}
		// 401
		if errors.As(err, &error.AuthUnauthorizedErr) {
			http.Error(rw, err.Error(), http.StatusUnauthorized)

			return
		}
		// 409
		if errors.As(err, &error.BllOrderAlreadyExistsAnotherErr) {
			http.Error(rw, err.Error(), http.StatusConflict)

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

	rw.WriteHeader(http.StatusAccepted)
}
