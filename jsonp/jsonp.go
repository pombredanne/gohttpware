package jsonp

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
)

func Handle(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		callback := r.URL.Query().Get("callback")
		if callback == "" {
			callback = r.URL.Query().Get("jsonp")
		}
		if callback == "" {
			h.ServeHTTP(w, r)
			return
		}

		wb := NewResponseBuffer(w)
		h.ServeHTTP(wb, r)

		if strings.Index(wb.Header().Get("Content-Type"), "/json") >= 0 {
			// Wrap the json data in the js callback function and set the
			// content type to javascript for the <script>
			wb.PreBody.Write([]byte(callback + "("))
			wb.Body.Write([]byte(")"))

			contentLength := wb.PreBody.Len() + wb.Body.Len()
			wb.Header().Set("Content-Type", "application/javascript")
			wb.Header().Set("Content-Length", strconv.Itoa(contentLength))
		}

		wb.Flush()
	}
	return http.HandlerFunc(fn)
}

type responseBuffer struct {
	Response http.ResponseWriter // the actual ResponseWriter to flush to
	Status   int                 // the HTTP response code from WriteHeader
	PreBody  *bytes.Buffer       // buffer to prepend to the content body
	Body     *bytes.Buffer       // the response content body
	Flushed  bool
}

func NewResponseBuffer(w http.ResponseWriter) *responseBuffer {
	return &responseBuffer{
		Response: w, Status: 200,
		PreBody: &bytes.Buffer{}, Body: &bytes.Buffer{},
	}
}

func (w *responseBuffer) Header() http.Header {
	return w.Response.Header() // use the actual response header
}

func (w *responseBuffer) Write(buf []byte) (int, error) {
	w.Body.Write(buf)
	return len(buf), nil
}

func (w *responseBuffer) WriteHeader(status int) {
	w.Status = status
}

func (w *responseBuffer) Flush() {
	if w.Flushed {
		return
	}

	w.Response.WriteHeader(w.Status)

	if w.PreBody.Len() > 0 {
		_, err := w.Response.Write(w.PreBody.Bytes())
		if err != nil {
			panic(err)
		}
		w.PreBody.Reset()
	}
	if w.Body.Len() > 0 {
		_, err := w.Response.Write(w.Body.Bytes())
		if err != nil {
			panic(err)
		}
		w.Body.Reset()
	}

	w.Flushed = true
}
