package logmw

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// LogResponseWriter is a custom type that extends http.ResponseWriter interface
// to capture and provide an easy access to http status code
type LogResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// Status is a easy way to retrieve the status code
func (w *LogResponseWriter) Status() int {
	return w.status
}

// Size provides the size of response object
func (w *LogResponseWriter) Size() int {
	return w.size
}

// Header returns the header to satisfy the http.ResponseWriter interface
func (w *LogResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write capture the size of the data written and satisfy the http.ResponseWriter interface
func (w *LogResponseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	written, err := w.ResponseWriter.Write(data)
	w.size += written
	return written, err
}

// WriteHeader capture the status code and satisfies the http.ResponseWriter interface
func (w *LogResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.status = code
}

func Middleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Try to get the real IP
		remoteAddr := r.RemoteAddr
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			remoteAddr = realIP
		}
		entry := log.WithFields(log.Fields{
			"request":    r.RequestURI,
			"method":     r.Method,
			"remote":     remoteAddr,
			"user-agent": r.UserAgent(),
			"referer":    r.Referer(),
		})
		res := &LogResponseWriter{ResponseWriter: w}
		h.ServeHTTP(res, r)
		latency := time.Since(start)
		var b strings.Builder
		switch {
		case latency < 500*time.Microsecond:
			fmt.Fprintf(&b, "%0.2f µsec", float64(latency.Nanoseconds())/float64(1000))
		case latency < 900*time.Millisecond:
			fmt.Fprintf(&b, "%0.2f msec", float64(latency.Nanoseconds())/float64(1000*1000))
		case latency < 100*time.Second:
			fmt.Fprintf(&b, "%0.2f sec", float64(latency.Nanoseconds())/float64(1000*1000*1000))
		}
		entry.WithFields(log.Fields{
			"status": res.Status(),
			"took":   b.String(),
			"size":   res.Size(),
		}).Info("completed handling request")
	}
	return http.HandlerFunc(fn)
}
