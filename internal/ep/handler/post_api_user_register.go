package handler

import (
	"errors"
	"net/http"

	errs "github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

func (cr *AppChiRouter) postApiUserRegister(rw http.ResponseWriter, r *http.Request) {
	jwtString, err := cr.authFacade.RegisterUser(r.Context(), r.Body)
	if err != nil {
		// 400
		// ..

		// 409
		if errors.As(err, &errs.ModelAlreadyExistsErr) {
			http.Error(rw, err.Error(), http.StatusConflict)

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
