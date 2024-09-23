package middleware

import (
	authService "Transactio/services/api-gateway/gRPC/proto"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const userKey = contextKey("userKey")
const requestId = contextKey("requestId")

var md metadata.MD
var requestUUID uuid.UUID

func MetaDataForGW(_ context.Context, _ *http.Request) metadata.MD {
	/*add metadata for grpc-gateway methods (login, signup)*/
	md = metadata.Pairs(string(requestId), requestUUID.String())
	return md
}

func LoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//middleware logic
			start := time.Now()

			requestUUID = uuid.New()
			ctx := context.WithValue(r.Context(), requestId, uuid.New())

			logRequest := fmt.Sprintf("%s Request %s from %s, reqId= %v", r.Method, r.RequestURI, r.RemoteAddr, requestUUID)
			logger.Info(logRequest)
			rsc := &ResponseStatusCode{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rsc, r.WithContext(ctx))

			logMsg := fmt.Sprintf("%s(%d %s) Response %s in %v, reqId= %v",
				r.Method, rsc.statusCode, http.StatusText(rsc.statusCode),
				r.RequestURI, time.Since(start), requestUUID)
			logger.Info(logMsg)
		})
	}
}

func AuthMiddleware(authClient authService.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//middleware logic

			if r.URL.Path == "/auth/login" || r.URL.Path == "/auth/signup" {
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := getToken(w, r)

			md = metadata.Pairs(string(requestId), requestUUID.String())
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			claims, err := authClient.ValidateJWT(ctx, &authService.JwtRequest{Token: tokenStr})
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				log.Println("Invalid token: ", err)
				return
			}

			ctx = context.WithValue(r.Context(), userKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CheckRole(roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, ok := r.Context().Value(userKey).(*authService.JwtResponse)
			if !ok || !containsRole(claims.Roles, roles) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

//-----
//UTILS
//-----

func containsRole(jwtRoles []string, acceptableRoles []string) bool {
	setAcceptableRoles := make(map[string]bool)
	for _, v := range acceptableRoles {
		setAcceptableRoles[v] = true
	}

	for _, jwtRole := range jwtRoles {
		if _, exist := setAcceptableRoles[jwtRole]; exist {
			return true
		}
	}
	return false
}

type ResponseStatusCode struct {
	http.ResponseWriter
	statusCode int
}

func (rst *ResponseStatusCode) WriteHeader(statusCode int) {
	rst.statusCode = statusCode
	rst.ResponseWriter.WriteHeader(statusCode)
}

func getToken(w http.ResponseWriter, r *http.Request) string {
	tokenHead := r.Header.Get("Authorization")
	if tokenHead == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return ""
	}

	tokenStr := strings.TrimPrefix(tokenHead, "Bearer ")
	return tokenStr
}
