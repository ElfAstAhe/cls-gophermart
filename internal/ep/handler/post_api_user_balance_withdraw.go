package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	_dto "github.com/ElfAstAhe/cls-gophermart/internal/ep/dto/v1"
	_utl "github.com/ElfAstAhe/cls-gophermart/internal/utils"
	_err "github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

func (cr *AppChiRouter) postApiUserBalanceWithdraw(rw http.ResponseWriter, r *http.Request) {
	cr.log.Debug("postApiUserBalanceWithdraw start")
	defer cr.log.Debug("postApiUserBalanceWithdraw finish")

	dec := json.NewDecoder(r.Body)
	defer _utl.CloseOnly(r.Body)
	var income _dto.WithdrawDto
	err := dec.Decode(&income)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	if err := cr.usersFacade.CreateWithdraw(r.Context(), &income); err != nil {
		// 401
		if errors.As(err, &_err.AuthUnauthorizedErr) {
			http.Error(rw, err.Error(), http.StatusUnauthorized)
		}
	}

	rw.WriteHeader(http.StatusOK)
}
