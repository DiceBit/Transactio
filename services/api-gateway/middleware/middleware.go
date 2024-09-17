package middleware

import (
	authService "Ecommers/services/api-gateway/gRPC/proto"
	"context"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const userKey = contextKey("userKey")

func AuthMiddleware(authClient authService.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//middleware logic
			log.Println("Middleware in Gateway")

			if r.URL.Path == "/auth/login" || r.URL.Path == "/auth/signup" {
				next.ServeHTTP(w, r)
			}

			tokenHead := r.Header.Get("Authorization")
			if tokenHead == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(tokenHead, "Bearer ")

			claims, err := authClient.ValidateJWT(context.Background(), &authService.JwtRequest{Token: tokenStr})
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				log.Println("Invalid token: ", err)
				return
			}

			ctx := context.WithValue(r.Context(), userKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthMiddleware2(next http.Handler, authClient authService.AuthServiceClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//middleware logic
		log.Println("Middleware in Gateway")

		if r.URL.Path == "/auth/login" || r.URL.Path == "/auth/signup" {
			next.ServeHTTP(w, r)
		}

		tokenHead := r.Header.Get("Authorization")
		if tokenHead == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(tokenHead, "Bearer ")

		claims, err := authClient.ValidateJWT(context.Background(), &authService.JwtRequest{Token: tokenStr})
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			log.Println("Invalid token: ", err)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
