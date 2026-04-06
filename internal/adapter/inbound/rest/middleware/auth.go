package middleware

import (
	"context"
	"net/http"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if token != "valid-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "userID", "12345"))

		next.ServeHTTP(w, r)
	})
}
