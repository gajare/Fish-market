package middleware

import (
	"net/http"

	"github.com/gajare/Fish-market/utils"
)

func AllowRoles(roles ...string) func(http.Handler) http.Handler {
	roleSet := map[string]struct{}{}
	for _, r := range roles {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, role, ok := GetUser(r)
			if !ok {
				utils.Error(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			if _, ok := roleSet[role]; !ok {
				utils.Error(w, http.StatusForbidden, "forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
