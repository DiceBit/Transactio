package middleware

import (
	authService "Transactio/internal/api-gateway/gRPC/authProto"
	"Transactio/internal/api-gateway/utils"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"time"
)

type contextKey string

const userKey = contextKey("userKey")
const requestId = contextKey("requestId")

var md metadata.MD
var requestUUID uuid.UUID

func MetaDataForGW(_ context.Context, _ *http.Request) metadata.MD {
	/*add metadata for grpc-api-gateway methods (login, signup)*/
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
			rsc := &utils.ResponseStatusCode{ResponseWriter: w, StatusCode: http.StatusOK}

			next.ServeHTTP(rsc, r.WithContext(ctx))

			code := rsc.StatusCode
			logMsg := fmt.Sprintf("%s(%d %s) Response %s in %v, reqId= %v",
				r.Method, rsc.StatusCode, http.StatusText(rsc.StatusCode),
				r.RequestURI, time.Since(start), requestUUID)

			if code <= 299 {
				logger.Info(logMsg)
			} else if code <= 599 {
				logger.Warn(logMsg)
			}

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

			tokenStr := utils.GetToken(w, r)

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
			if !ok || !utils.ContainsRole(claims.Roles, roles) {

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
