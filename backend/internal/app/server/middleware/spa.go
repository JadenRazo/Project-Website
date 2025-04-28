package middleware

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// SPAConfig contains configuration for the SPA middleware
type SPAConfig struct {
	// StaticPath is the path to the built static files
	StaticPath string

	// IndexPath is the path to the index.html file
	IndexPath string

	// ApiPrefix is the prefix for API routes that should bypass the SPA handling
	ApiPrefix string

	// AssetExtensions is a list of file extensions that should be treated as static assets
	AssetExtensions []string
}

// DefaultSPAConfig returns the default SPA configuration
func DefaultSPAConfig() SPAConfig {
	return SPAConfig{
		StaticPath:      "./frontend/build",
		IndexPath:       "index.html",
		ApiPrefix:       "/api",
		AssetExtensions: []string{".js", ".css", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".woff", ".woff2", ".ttf", ".eot"},
	}
}

// SPA middleware serves a Single Page Application and handles client-side routing
// It serves static files directly, passes API requests through, and serves the index.html
// for all other requests to enable client-side routing
func SPA(config SPAConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// If no config provided, use defaults
		if config.StaticPath == "" {
			config = DefaultSPAConfig()
		}

		// Make sure paths are properly formatted
		staticPath := filepath.Clean(config.StaticPath)
		indexPath := filepath.Join(staticPath, config.IndexPath)

		// Create the file server for static assets
		fileServer := http.FileServer(http.Dir(staticPath))

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the path
			path := r.URL.Path

			// If it's an API request, skip SPA handling and pass to the next handler
			if strings.HasPrefix(path, config.ApiPrefix) {
				next.ServeHTTP(w, r)
				return
			}

			// Check if this is a request for a static asset
			isAsset := false
			for _, ext := range config.AssetExtensions {
				if strings.HasSuffix(path, ext) {
					isAsset = true
					break
				}
			}

			// If it's a static asset, try to serve it directly
			if isAsset {
				// Construct the file path
				filePath := filepath.Join(staticPath, path)

				// Check if the file exists
				_, err := os.Stat(filePath)
				if err == nil {
					// File exists, serve it
					fileServer.ServeHTTP(w, r)
					return
				}
			}

			// For all other requests, serve the index.html to enable client-side routing
			// First, check if index.html exists
			_, err := os.Stat(indexPath)
			if os.IsNotExist(err) {
				// If index.html doesn't exist, pass to the next handler
				next.ServeHTTP(w, r)
				return
			}

			// Serve the index.html file
			http.ServeFile(w, r, indexPath)
		})
	}
}

// SPAHandler returns a simple http.Handler that serves a Single Page Application
// This is a simpler version that can be used directly as a handler
func SPAHandler(staticPath string, indexPath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get absolute path to requested file
		path := filepath.Join(staticPath, r.URL.Path)

		// Check if file exists
		_, err := os.Stat(path)

		// If file exists and is not a directory, serve it directly
		if err == nil {
			fileInfo, _ := os.Stat(path)
			if !fileInfo.IsDir() {
				http.ServeFile(w, r, path)
				return
			}
		}

		// Otherwise serve the index.html
		http.ServeFile(w, r, filepath.Join(staticPath, indexPath))
	})
}
