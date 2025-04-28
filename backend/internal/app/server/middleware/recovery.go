package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/JadenRazo/Project-Website/backend/internal/common/logger"
)

// RecoveryOptions configures the recovery middleware
type RecoveryOptions struct {
	// EnableStackTrace determines if stack traces should be included in the response
	EnableStackTrace bool

	// ResponseWriter is a custom function to write the response when a panic occurs
	ResponseWriter func(w http.ResponseWriter, r *http.Request, err interface{}, stack []byte)
}

// DefaultRecoveryOptions returns the default recovery options
func DefaultRecoveryOptions() *RecoveryOptions {
	return &RecoveryOptions{
		EnableStackTrace: false,
		ResponseWriter:   defaultRecoveryResponse,
	}
}

// defaultRecoveryResponse is the default response writer for recovered panics
func defaultRecoveryResponse(w http.ResponseWriter, r *http.Request, err interface{}, stack []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	// Don't expose internal error details in production
	w.Write([]byte(`{"error":"Internal Server Error"}`))
}

// Recovery middleware recovers from panics and logs the error
func Recovery(opts *RecoveryOptions) func(next http.Handler) http.Handler {
	if opts == nil {
		opts = DefaultRecoveryOptions()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// Capture stack trace
					stack := debug.Stack()

					// Convert panic value to string
					var err error
					switch recVal := rec.(type) {
					case string:
						err = errors.New(recVal)
					case error:
						err = recVal
					default:
						err = fmt.Errorf("%v", rec)
					}

					// Get request details for logging
					requestID := r.Header.Get("X-Request-ID")
					if requestID == "" {
						requestID = "unknown"
					}

					// Log the panic
					logger.Error("Panic recovered",
						"request_id", requestID,
						"error", err.Error(),
						"stack", string(stack),
						"path", r.URL.Path,
						"method", r.Method,
					)

					// Write response
					opts.ResponseWriter(w, r, rec, stack)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// JSONRecoveryResponse writes a JSON response for recovered panics
// This is useful for API servers where all responses should be JSON
func JSONRecoveryResponse(w http.ResponseWriter, r *http.Request, err interface{}, stack []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	// In development mode, you might want to include more details
	isDev := r.Header.Get("X-Environment") == "development"

	if isDev {
		errMsg := fmt.Sprintf("%v", err)
		w.Write([]byte(fmt.Sprintf(`{"error":"Internal Server Error","message":"%s"}`, errMsg)))
	} else {
		w.Write([]byte(`{"error":"Internal Server Error"}`))
	}
}
