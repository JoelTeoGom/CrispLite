package middleware

import (
	"crisplite/internal/port/outbound"
	"fmt"
	"net/http"
	"time"
)

func Logger(logger outbound.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(r.Context(), fmt.Sprintf("%s | %s | %s | %s",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.RemoteAddr,
			r.URL.Path,
		))
		next.ServeHTTP(w, r)
	})
}
