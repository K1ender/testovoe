package utils

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// gzipResponseWriter wraps http.ResponseWriter to write gzipped content.
type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.writer.Write(b)
}

func (g *gzipResponseWriter) WriteHeader(statusCode int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(statusCode)
}

// GzipMiddleware compresses HTTP responses for clients that support gzip.
// It skips compression if the client does not advertise gzip support or if
// the response is already encoded.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Do not double-encode if content is pre-encoded
		if enc := w.Header().Get("Content-Encoding"); enc != "" && enc != "identity" {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		grw := &gzipResponseWriter{ResponseWriter: w, writer: gz}
		next.ServeHTTP(grw, r)
	})
}

// GzipWriter wraps an io.Writer with gzip. Useful for tests or utilities.
func GzipWriter(w io.Writer) *gzip.Writer {
	return gzip.NewWriter(w)
}

