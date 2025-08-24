package middleware

import (
	"net/http"
	"time"

	"github.com/gajare/Fish-market/logger"
	"github.com/gajare/Fish-market/pkg/ctxkeys"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(code int) { w.status = code; w.ResponseWriter.WriteHeader(code) }
func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w}

		next.ServeHTTP(sw, r)

		reqID, _ := r.Context().Value(ctxkeys.ReqID).(string)
		uid, _ := r.Context().Value(ctxkeys.UserID).(uint)
		role, _ := r.Context().Value(ctxkeys.Role).(string)

		logger.With(map[string]any{
			"level":      "info",
			"req_id":     reqID,
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     sw.status,
			"bytes":      sw.length,
			"latency_ms": time.Since(start).Milliseconds(),
			"ip":         r.RemoteAddr,
			"user_agent": r.UserAgent(),
			"user_id":    uid,
			"role":       role,
		}).Info("http_request")
	})
}
