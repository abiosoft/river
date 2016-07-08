package river

import "net/http"

// customrw is a custom ResponseWriter.
type customrw struct {
	header        func() http.Header
	write         func([]byte) (int, error)
	writeHeader   func(int)
	headerWritten bool
}

func (r *customrw) Header() http.Header {
	return r.header()
}

func (r *customrw) Write(b []byte) (int, error) {
	if !r.headerWritten {
		r.WriteHeader(http.StatusOK)
	}
	return r.write(b)
}

func (r *customrw) WriteHeader(status int) {
	r.headerWritten = true
	r.writeHeader(status)
}

// noWriteRW returns a wrapper of w with invalidated Write and WriteHeader.
func noWriteRW(w http.ResponseWriter) http.ResponseWriter {
	rw := &customrw{
		header: func() http.Header {
			return w.Header()
		},
		write: func(b []byte) (int, error) {
			return len(b), nil
		},
		writeHeader: func(status int) {},
	}
	return rw
}

// beforeWriteRW returns a wrapper of w that executes f before writing.
func beforeWriteRW(w http.ResponseWriter, f func()) http.ResponseWriter {
	rw := &customrw{
		header: func() http.Header {
			return w.Header()
		},
		write: func(b []byte) (int, error) {
			return w.Write(b)
		},
		writeHeader: func(status int) {
			f()
			w.WriteHeader(status)
		},
	}
	return rw
}

// staticStatusRW returns a wrapper of w that sets status as status code.
func staticStatusRW(w http.ResponseWriter, status int) http.ResponseWriter {
	rw := &customrw{
		header: func() http.Header {
			return w.Header()
		},
		write: func(b []byte) (int, error) {
			return w.Write(b)
		},
		writeHeader: func(_ int) {
			w.WriteHeader(status)
		},
	}
	return rw
}
