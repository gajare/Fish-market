package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gajare/Fish-market/pkg/ctxkeys"
	"github.com/gajare/Fish-market/utils"
)

type ctxKey string

const (
	ctxUserID   ctxKey = "uid"
	ctxUserRole ctxKey = "role"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer") {
			utils.Error(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		raw := strings.TrimPrefix(h, "Bearer ")
		claims, err := utils.ParseJWT(raw)
		if err != nil {
			utils.Error(w, http.StatusUnauthorized, "invalid token")
			return
		}
		ctx := context.WithValue(r.Context(), ctxkeys.UserID, claims.UserID)
		ctx = context.WithValue(ctx, ctxkeys.Role, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(r *http.Request) (uint, string, bool) {
	id, ok1 := r.Context().Value(ctxUserID).(uint)
	role, ok2 := r.Context().Value(ctxUserRole).(string)
	return id, role, ok1 && ok2
}
