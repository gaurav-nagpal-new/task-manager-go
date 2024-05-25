package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"task-manager/constants"
	"task-manager/jwtauth"
	"task-manager/utils"

	"github.com/golang-jwt/jwt"
)

func VerifyJWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := tokenParts[1]
		claims := &jwtauth.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv(constants.JwtTokenKey)), nil
			})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), constants.UserCollectionName, utils.GetTaskCollectionName(claims.Email))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
