package server

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

type contextKey string

const requestId = contextKey("requestId")

func LogInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.DataLoss, "Missing metadata")
		}

		requestUUID := md.Get(string(requestId))

		infoMsg := fmt.Sprintf("UnaryCall(%s) Request %s",
			requestUUID, info.FullMethod)
		logger.Info(infoMsg)

		resp, err := handler(ctx, req)
		if err != nil {
			errMsg := fmt.Sprintf("UnaryCall(%s) %s in %v Error: %v",
				requestUUID, info.FullMethod, time.Since(start), err)
			logger.Error(errMsg)
		} else {
			infoMsg = fmt.Sprintf("UnaryCall(%s) Response %s in %v",
				requestUUID, info.FullMethod, time.Since(start))
			logger.Info(infoMsg)
		}

		return resp, nil
	}
}
