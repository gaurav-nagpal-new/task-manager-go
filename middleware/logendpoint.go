package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func LogEndPoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zap.L().Debug("API requested", zap.String("URL", r.URL.String()), zap.String("Method", r.Method))
		next.ServeHTTP(w, r)
	})
}
