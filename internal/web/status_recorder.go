package web

import "net/http"

// StatusRecorder embeds an [http.ResponseWriter] which remembers the status
// code being generated, allowing client to retroactively query the status code.
type StatusRecorder struct {
	http.ResponseWriter
	code int
}

func (r *StatusRecorder) WriteHeader(code int) {
	r.ResponseWriter.WriteHeader(code)
	r.code = code
}

// Unwrap implements the unexported rwUnwrapper interface. This is necessary for
// [http.ResponseController] to get the underlying ResponseWriter, e.g. to
// query for cabilities like [http.Flusher].
func (r *StatusRecorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }

func (r *StatusRecorder) Code() int {
	if r.code == 0 {
		return 200
	}
	return r.code
}
