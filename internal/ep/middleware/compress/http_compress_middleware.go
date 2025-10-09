package compress

import (
	"io"
	"net/http"

	"github.com/andybalholm/brotli"
	"github.com/go-chi/chi/v5/middleware"
)

func CustomCompress(level int, allowedTypes ...string) func(next http.Handler) http.Handler {
	compressor := middleware.NewCompressor(level, allowedTypes...)
	// add brotli
	compressor.SetEncoder(encodingBrotli, encoderBrotli)

	return compressor.Handler
}

func encoderBrotli(w io.Writer, level int) io.Writer {
	br := brotli.NewWriterLevel(w, level)

	return br
}
