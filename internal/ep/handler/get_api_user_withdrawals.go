package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

func (cr *AppChiRouter) getAPIUserWithdrawals(rw http.ResponseWriter, r *http.Request) {
	cr.log.Debug("getAPIUserWithdrawals start")
	defer cr.log.Debug("getAPIUserWithdrawals finish")

	dtoList, err := cr.usersFacade.ListWithdrawals(r.Context())
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

	if len(dtoList) == 0 {
		// 204
		rw.WriteHeader(http.StatusNoContent)
	} else {
		// 200
		rw.WriteHeader(http.StatusOK)
	}
	rw.Header().Add("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	if err := enc.Encode(dtoList); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
