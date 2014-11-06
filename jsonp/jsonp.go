package jsonp

import (
	"bytes"
	"log"
	"net/http"
)

func Handle(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// check the request.. is this jsonp shiz.. if so, let's
		// make sure we support the callback and return on the way out..
		log.Println("JSONP, BEFORE....... is this json...? check callback= and headers..")

		callback := r.URL.Query().Get("callback")
		if callback == "" {
			h.ServeHTTP(w, r)
			return
		}

		wb := NewResponseBuffer(w)
		wb.Body.Write([]byte(callback + "("))
		h.ServeHTTP(wb, r)
		wb.Body.Write([]byte(")"))
		wb.Flush()
	}
	return http.HandlerFunc(fn)
}

// Response buffer, based on httptest.ResponseRecorder
type ResponseBuffer struct {
	RW        http.ResponseWriter // The actual ResponseWriter to flush to
	Code      int                 // the HTTP response code from WriteHeader
	HeaderMap http.Header         // the HTTP response headers
	Body      *bytes.Buffer       // if non-nil, the bytes.Buffer to append written data to
	Flushed   bool

	wroteHeader bool
}

func NewResponseBuffer(w http.ResponseWriter) *ResponseBuffer {
	return &ResponseBuffer{
		RW: w, HeaderMap: make(http.Header), Body: &bytes.Buffer{},
	}
}

func (w *ResponseBuffer) Header() http.Header {
	return w.HeaderMap
	// if m == nil {
	// 	m = make(http.Header)
	// 	w.HeaderMap = m
	// }
	// return m
}

func (w *ResponseBuffer) Write(buf []byte) (int, error) {
	// if !w.wroteHeader {
	// 	w.WriteHeader(200)
	// }
	if w.Body == nil {
		w.Body = bytes.NewBuffer(buf)
	} else {
		w.Body.Write(buf)
	}
	return len(buf), nil
}

func (w *ResponseBuffer) WriteHeader(code int) {
	if !w.wroteHeader {
		w.Code = code
	}
	w.wroteHeader = true
}

func (w *ResponseBuffer) Flush() {
	if !w.Flushed {
		if !w.wroteHeader {
			w.WriteHeader(200)
		}
		w.RW.WriteHeader(w.Code)
	}

	n, err := w.RW.Write(w.Body.Bytes())
	_ = n
	_ = err

	w.Body.Reset()
	w.Flushed = true
}
