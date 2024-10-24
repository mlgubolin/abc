package http

import (
	"application"
	"encoding/json"
	"net/http"
	"runtime"
)

// Error prints & optionally logs an error message.
func (s *Server) Error(w http.ResponseWriter, r *http.Request, err error) {
	// Extract error code & message.
	code, message := application.ErrorCode(err), application.ErrorMessage(err)

	// Log the error with additional details.
	clientIP := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")
	referer := r.Header.Get("Referer")

	// Capture the file and line number where the error occurred.
	_, file, line, _ := runtime.Caller(1)

	s.logger.Printf("http error: method=%s path=%s client_ip=%s user_agent=%s referer=%s error=%v code=%s message=%s file=%s line=%d",
		r.Method, r.URL.Path, clientIP, userAgent, referer, err, code, message, file, line)

	// Log & report internal errors.
	if code == application.ErrInternal {
		application.ReportError(r.Context(), err, r)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(ErrorStatusCode(code))

	json.NewEncoder(w).Encode(&ErrorResponse{Error: message})
}

// ErrorResponse represents a JSON structure for error output.
type ErrorResponse struct {
	Error string `json:"error"`
}

// lookup of application error codes to HTTP status codes.
var codes = map[string]int{
	application.ErrConflict:       http.StatusConflict,
	application.ErrInvalid:        http.StatusBadRequest,
	application.ErrNotFound:       http.StatusNotFound,
	application.ErrNotImplemented: http.StatusNotImplemented,
	application.ErrUnauthorized:   http.StatusUnauthorized,
	application.ErrNotAllowed:     http.StatusMethodNotAllowed,
	application.ErrInternal:       http.StatusInternalServerError,
}

// ErrorStatusCode returns the associated HTTP status code for a Abacateiro error code.
func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

// FromErrorStatusCode returns the associated Abacateiro code for an HTTP status code.
// func FromErrorStatusCode(code int) string {
// 	for k, v := range codes {
// 		if v == code {
// 			return k
// 		}
// 	}
// 	return application.ErrInternal
// }

// func (s *Server) notFoundResponse(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "/ui", http.StatusSeeOther)
// }

// func (s *Server) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {

// 	err := application.Errorf(
// 		application.ErrNotAllowed,
// 		"the %s method is not supported for this resource",
// 		r.Method)

// 	s.Error(w, r, err)
// }
