package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

func (cr *AppChiRouter) getApiUserBalance(rw http.ResponseWriter, r *http.Request) {
	cr.log.Debug("getApiUserBalance start")
	defer cr.log.Debug("getApiUserBalance finish")

	// userID in request context
	dto, err := cr.usersFacade.GetBalance(r.Context())
	// error check
	if err != nil {
		// 401
		if errors.As(err, &error.AuthUnauthorizedErr) {
			http.Error(rw, err.Error(), http.StatusUnauthorized)

			return
		}

		// 500
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	if err := enc.Encode(dto); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
