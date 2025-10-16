package handler

import (
	"errors"
	"io"
	"net/http"

	_utl "github.com/ElfAstAhe/cls-gophermart/internal/utils"
	_err "github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

func (cr *AppChiRouter) postApiUserOrders(rw http.ResponseWriter, r *http.Request) {
	cr.log.Debug("postApiUserOrders start")
	defer cr.log.Debug("postApiUserOrders finish")

	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(rw, "Content-Type not supported, supported text/plain", http.StatusBadRequest)

		return
	}
	bytes, err := io.ReadAll(r.Body)
	defer _utl.CloseOnly(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)

		return
	}

	// ToDo: make a table ? think
	if err := cr.usersFacade.CreateOrder(r.Context(), bytes); err != nil {
		// 200
		if errors.As(err, &_err.BllOrderAlreadyExistsCurrentErr) {
			rw.WriteHeader(http.StatusOK)

			return
		}
		// 400
		if errors.As(err, &_err.ModelValidationErr) {
			http.Error(rw, err.Error(), http.StatusBadRequest)

			return
		}
		// 401
		if errors.As(err, &_err.AuthUnauthorizedErr) {
			http.Error(rw, err.Error(), http.StatusUnauthorized)

			return
		}
		// 409
		if errors.As(err, &_err.BllOrderAlreadyExistsAnotherErr) {
			http.Error(rw, err.Error(), http.StatusConflict)

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

	rw.WriteHeader(http.StatusAccepted)
}
