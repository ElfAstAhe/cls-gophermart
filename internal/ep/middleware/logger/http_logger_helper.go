package logger

import (
	"net/http"
)

// little helper for log http response info
type responseInfo struct {
	StatusCode int
	Size       int64
}

func emptyResponseInfo() *responseInfo {
	return newResponseInfo(-1, -1)
}

func newResponseInfo(statusCode int, size int64) *responseInfo {
	return &responseInfo{StatusCode: statusCode, Size: size}
}

// ResponseLoggerWriter http response logging
type ResponseLoggerWriter struct {
	http.ResponseWriter
	info *responseInfo
}

func NewResponseLoggerWriter(rw http.ResponseWriter) *ResponseLoggerWriter {
	return &ResponseLoggerWriter{
		ResponseWriter: rw,
		info:           emptyResponseInfo(),
	}
}

// http.ResponseWriter =========================

func (rlw *ResponseLoggerWriter) WriteHeader(statusCode int) {
	rlw.ResponseWriter.WriteHeader(statusCode)
	rlw.info.StatusCode = statusCode
}

func (rlw *ResponseLoggerWriter) Write(b []byte) (int, error) {
	size, err := rlw.ResponseWriter.Write(b)
	rlw.info.Size += int64(size)

	return size, err
}

// =============================================
