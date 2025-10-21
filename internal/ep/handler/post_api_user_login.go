package handler

import (
	"errors"
	"net/http"

	errs "github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

func (cr *AppChiRouter) postApiUserLogin(rw http.ResponseWriter, r *http.Request) {
	jwtString, err := cr.authFacade.LoginUser(r.Context(), r.Body)
	if err != nil {
		// 400
		// ..

		// 401
		if errors.As(err, &errs.AuthUnauthorizedErr) {
			http.Error(rw, err.Error(), http.StatusUnauthorized)

			return
		}

		// 500
		http.Error(rw, err.Error(), http.StatusInternalServerError)

		return
	}

	cr.setJWTAuthCookie(jwtString, rw)

	// 200
	rw.WriteHeader(http.StatusOK)
}
