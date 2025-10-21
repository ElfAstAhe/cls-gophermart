package jwt

import (
	"context"
	"net/http"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/service/auth"
)

type HTTPJWTMiddleware struct {
	log logger.Logger
}

func NewHTTPJWTMiddleware(logger logger.Logger) *HTTPJWTMiddleware {
	return &HTTPJWTMiddleware{
		log: logger.GetLogger("jwt_middleware"),
	}
}

func (m *HTTPJWTMiddleware) UserInfoRetrieve(next http.Handler) http.Handler {
	m.log.Debug("User Info Retrieve start")
	defer m.log.Debug("User Info Retrieve finish")

	fn := func(w http.ResponseWriter, r *http.Request) {
		userInfo, err := auth.UserInfoFromRequestJWT(r)
		m.log.Debugf("User Info Retrieve userInfo [%v] err [%v]", userInfo, err)
		req := r
		if err != nil {
			m.log.Warnf("no jwt info with error: [%v]", err)
		} else {
			req = r.WithContext(context.WithValue(r.Context(), auth.ContextUserInfo, userInfo))
		}

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}
