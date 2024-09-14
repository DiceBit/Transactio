package web

import (
	"context"
	"log"
	"net/http"
	"strings"
	"user-service/pkg/utils"
)

type contextKey string

const userKey = contextKey("userKey")

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//middleware logic
		tokenHead := r.Header.Get("Authorization")
		if tokenHead == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(tokenHead, "Bearer ")
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			log.Println("Invalid token: ", err)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleMiddleware(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//middleware logic

			claims, ok := r.Context().Value(userKey).(*utils.Claims)
			if !ok || !contains(claims.Roles, role) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func contains(roles []string, role string) bool {
	for _, v := range roles {
		if v == role {
			return true
		}
	}
	return false
}
