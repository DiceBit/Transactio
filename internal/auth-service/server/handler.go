package server

import (
	"Transactio/internal/auth-service/db"
	pb "Transactio/internal/auth-service/gRPC/proto"
	userUtils "Transactio/internal/auth-service/utils"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
	grpcServer *grpc.Server
	db         *pgxpool.Pool

	logger *zap.Logger

	authName string
	authAddr string
}

func (authServ *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	logger := authServ.logger

	usr, err := pgxDb.UsrByEmail(ctx, authServ.db, req.Email)
	if err != nil {
		logger.Error(
			"Error with getting usr by email",
			zap.String("emailRequest", req.Email),
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal,
			"Error with getting usr by email(%s). %v", req.Email, err)
	}

	if ok := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(req.Password)); ok != nil {
		logger.Warn(
			"Invalid password",
			zap.String("email", usr.Email),
		)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid password for %s", usr.Email)
	}

	token, err := userUtils.GenerateJWT(usr.Email, usr.Role)
	if err != nil {
		logger.Error(
			"Error generating token",
			zap.String("emailForJwtGen", usr.Email),
			zap.Strings("roleForJwtGen", usr.Role),
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "Error generating token, %s %s",
			usr.Email, usr.Role)
	}

	return &pb.LoginResponse{Token: token}, nil
}

// registration
func (authServ *AuthServiceServer) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	logger := authServ.logger
	exist, err := pgxDb.CheckIfExistUsr(ctx, authServ.db, req)
	if err != nil {
		logger.Error(
			"Error with checking usr in DB",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		return nil, status.Errorf(codes.Internal, "Error with checking usr(%s) in DB: %v", req.Username, err)
	}
	if exist {
		logger.Warn(
			"User already exists",
			zap.String("username", req.Username),
		)
		return nil, status.Errorf(codes.AlreadyExists, "User(%s) already exists", req.Username)
	} else {
		err = pgxDb.AddUser(ctx, authServ.db, req)
		if err != nil {
			logger.Error(
				"Error with adding usr",
				zap.String("username", req.Username),
				zap.Error(err),
			)
			return nil, status.Errorf(codes.Internal, "Error with adding usr(%s). %v", req.Username, err)
		}
		logger.Info(
			"User added",
			zap.String("username", req.Username),
		)
	}

	usr, err := pgxDb.UsrByEmail(ctx, authServ.db, req.Email)
	if err != nil {
		logger.Error("Error with getting usr by email",
			zap.String("email", req.Email))
		return nil, status.Errorf(codes.Internal, "Error with getting usr by email(%s)", req.Email)
	}

	token, err := userUtils.GenerateJWT(usr.Email, usr.Role)
	if err != nil {
		logger.Error("Error generating token",
			zap.String("email", usr.Email),
			zap.Strings("roles", usr.Role))
		return nil, status.Errorf(codes.Internal, "Error generating token, %s %s",
			usr.Email, usr.Role)
	}

	return &pb.SignUpResponse{Token: token}, nil
}

func (authServ *AuthServiceServer) ValidateJWT(_ context.Context, req *pb.JwtRequest) (*pb.JwtResponse, error) {
	claims, err := userUtils.ValidateJWT(req.Token)
	if err != nil {
		authServ.logger.Error("Error with validation jwt", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "Error with validation jwt. %v", err)
	}

	exp := claims.ExpiresAt.Unix()
	iat := claims.IssuedAt.Unix()
	nbr := claims.NotBefore.Unix()

	return &pb.JwtResponse{Email: claims.Email, Roles: claims.Roles,
		Exp: exp, Iat: iat, Nbr: nbr}, nil
}
