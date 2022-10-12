package CADDY_FILE_PROXY

import (
	"bytes"
	"io"
	"net/http"
)

type ResponseWrapper struct {
	w    http.ResponseWriter
	buf  bytes.Buffer
	code int
}

func (rw *ResponseWrapper) Header() http.Header {
	return rw.w.Header()
}

func (rw *ResponseWrapper) WriteHeader(statusCode int) {
	rw.code = statusCode
}

func (rw *ResponseWrapper) Write(data []byte) (int, error) {
	return rw.buf.Write(data)
}
func (rw *ResponseWrapper) Done() (int64, error) {
	if rw.code > 0 {
		rw.w.WriteHeader(rw.code)
	}
	return io.Copy(rw.w, &rw.buf)
}
