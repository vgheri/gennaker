package route

import (
	"net/http"
)

type responseWriterWrapper struct {
	statusCode int
	http.ResponseWriter
}

func newResponseWriterWrapper(res http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{200, res}
}

func (wrapper *responseWriterWrapper) Status() int {
	return wrapper.statusCode
}

// Satisfy the http.ResponseWriter interface
func (wrapper *responseWriterWrapper) Header() http.Header {
	return wrapper.ResponseWriter.Header()
}

func (wrapper *responseWriterWrapper) Write(content []byte) (int, error) {
	return wrapper.ResponseWriter.Write(content)
}

func (wrapper *responseWriterWrapper) WriteHeader(code int) {
	wrapper.statusCode = code
	wrapper.ResponseWriter.WriteHeader(code)
}
