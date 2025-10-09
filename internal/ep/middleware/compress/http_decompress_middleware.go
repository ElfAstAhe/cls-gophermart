package compress

import (
	"net/http"
	"strings"

	_utl "github.com/ElfAstAhe/cls-gophermart/internal/utils"
)

func CustomDecompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoding := strings.ToLower(r.Header.Get("Content-Encoding"))
		if encoding != "" {
			dr, err := NewCustomReader(r.Body, encoding)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}
			defer _utl.CloseOnly(dr)

			r.Body = dr
		}

		// no decompression
		next.ServeHTTP(w, r)
	})
}
