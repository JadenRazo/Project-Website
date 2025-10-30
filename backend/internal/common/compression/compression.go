package compression

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/JadenRazo/Project-Website/backend/internal/app/config"
)

var (
	gzipPool = sync.Pool{
		New: func() interface{} {
			return gzip.NewWriter(nil)
		},
	}
	zlibPool = sync.Pool{
		New: func() interface{} {
			return zlib.NewWriter(nil)
		},
	}
)

// CompressionMiddleware creates a middleware that compresses HTTP responses
func CompressionMiddleware(cfg *config.CompressionConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip if compression is disabled
			if !cfg.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Check if client supports compression
			acceptEncoding := r.Header.Get("Accept-Encoding")
			if acceptEncoding == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Determine compression method
			var writer io.Writer
			var compressionType string

			if strings.Contains(acceptEncoding, "gzip") {
				gz := gzipPool.Get().(*gzip.Writer)
				gz.Reset(w)
				writer = gz
				compressionType = "gzip"
			} else if strings.Contains(acceptEncoding, "deflate") {
				zl := zlibPool.Get().(*zlib.Writer)
				zl.Reset(w)
				writer = zl
				compressionType = "deflate"
			} else {
				next.ServeHTTP(w, r)
				return
			}

			// Create compressed response writer
			compressedWriter := &compressedResponseWriter{
				ResponseWriter: w,
				writer:         writer,
				compression:    compressionType,
			}

			// Set compression headers
			w.Header().Set("Content-Encoding", compressionType)
			w.Header().Set("Vary", "Accept-Encoding")

			// Handle the request
			next.ServeHTTP(compressedWriter, r)

			// Close and return writer to pool
			if gz, ok := writer.(*gzip.Writer); ok {
				gz.Close()
				gzipPool.Put(gz)
			} else if zl, ok := writer.(*zlib.Writer); ok {
				zl.Close()
				zlibPool.Put(zl)
			}
		})
	}
}

type compressedResponseWriter struct {
	http.ResponseWriter
	writer      io.Writer
	compression string
}

func (c *compressedResponseWriter) Write(b []byte) (int, error) {
	return c.writer.Write(b)
}

func (c *compressedResponseWriter) WriteHeader(statusCode int) {
	c.ResponseWriter.WriteHeader(statusCode)
}

func (c *compressedResponseWriter) Flush() {
	if f, ok := c.writer.(http.Flusher); ok {
		f.Flush()
	}
}
