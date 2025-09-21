package librest

import (
	"encoding/json"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
	err         error
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default http status OK
	}
}

// WrapResponseWriter wrap original http.ResponseWriter with custom ResponseWriter
func WrapResponseWriter(w http.ResponseWriter) *ResponseWriter {
	var rw *ResponseWriter
	rw, ok := w.(*ResponseWriter)
	if !ok {
		rw = NewResponseWriter(w)
	}

	return rw
}

func (rw *ResponseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.statusCode = code
		rw.ResponseWriter.WriteHeader(code)
		rw.wroteHeader = true
	}
}

func (rw *ResponseWriter) SetError(err error) {
	rw.err = err
}

func (rw *ResponseWriter) GetError() error {
	return rw.err
}

func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

func WriteHTTPResponse(w http.ResponseWriter, body any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}
