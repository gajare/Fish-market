package middleware

import (
	"context"
	"net/http"

	"github.com/gajare/Fish-market/pkg/ctxkeys"
	"github.com/google/uuid"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-Id")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		w.Header().Set("X-Request-Id", reqID)
		ctx := context.WithValue(r.Context(), ctxkeys.ReqID, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
