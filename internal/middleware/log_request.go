package middleware

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)


type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func (w *statusWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *statusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return h.Hijack()
}

func (w *statusWriter) Push(target string, opts *http.PushOptions) error {
	p, ok := w.ResponseWriter.(http.Pusher)
	if !ok {
		return http.ErrNotSupported
	}
	return p.Push(target, opts)
}

func (w *statusWriter) ReadFrom(r io.Reader) (int64, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}

	rf, ok := w.ResponseWriter.(io.ReaderFrom)
	if !ok {
		return io.Copy(w, r)
	}

	n, err := rf.ReadFrom(r)
	w.bytes += int(n)
	return n, err
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (m *Middleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ShouldInstrument(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		sw := &statusWriter{ResponseWriter: w}

		next.ServeHTTP(sw, r)

		status := sw.status
		if status == 0 {
			status = http.StatusOK
		}

		log := LoggerFrom(r.Context())
		log.Info("http request",
			zap.String("request_id", RequestIDFrom(r.Context())),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.Int("status", status),
			zap.Int("bytes", sw.bytes),
			zap.Duration("duration", time.Since(start)),
			zap.String("remote_ip", clientIP(r)),
			zap.String("user_agent", r.UserAgent()),
		)
	})
}