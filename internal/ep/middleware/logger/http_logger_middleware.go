package logger

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
)

type HTTPLoggerMiddleware struct {
	log logger.Logger
}

func NewHTTPLoggerMiddleware(logger logger.Logger) *HTTPLoggerMiddleware {
	return &HTTPLoggerMiddleware{
		log: logger,
	}
}

func (hl *HTTPLoggerMiddleware) CustomInfoHTTPLogger(nextHandler http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := NewResponseLoggerWriter(rw)

		nextHandler.ServeHTTP(lrw, r)

		duration := time.Since(start)

		hl.log.Infof("uri [%s] method [%s] duration [%v]ms status [%v] size [%v]",
			r.RequestURI, r.Method, strconv.FormatInt(duration.Milliseconds(), 10),
			lrw.info.StatusCode, lrw.info.Size)
	}

	return http.HandlerFunc(fn)
}
