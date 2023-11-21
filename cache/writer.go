package cache

import (
	"net/http"
)

// wrapper for response writer that caches the response
type Writer struct {
	writer http.ResponseWriter
	response response
	resource string
}

// interface implementation check
// unused variable of type http.ResponseWriter 
// that is equal to nil being casted into a pointer
// to Writer struct
var (
	_ http.ResponseWriter = (*Writer)(nil)
)

// return new cache writer
func NewWriter(w http.ResponseWriter, r *http.Request) *Writer {
	return &Writer{
		writer: w,
		resource: MakeResource(r),
		response: response{
			header: http.Header{},
		},
	}
}

// returns the response headers
func (w *Writer) Header() http.Header {
	return w.response.header
}

func (w *Writer) WriteHeader(code int) {
	copyHeader(w.response.header, w.writer.Header())
	w.response.code = code
	w.writer.WriteHeader(code)
}

func (w *Writer) Write(b []byte) (int, error) {
	w.response.body = make([]byte, len(b))

	for k, v := range b {
		w.response.body[k] = v
	}

	copyHeader(w.Header(), w.writer.Header())
	set(w.resource, &w.response)
	return w.writer.Write(b)
}
