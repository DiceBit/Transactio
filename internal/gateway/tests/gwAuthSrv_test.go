package tests

import (
	pb "Transactio/internal/gateway/gRPC/proto"
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"testing"

	_ "github.com/stretchr/testify/assert"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func TestMain(m *testing.M) {
	lis = bufconn.Listen(bufSize)
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, &mockAuthService{})

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	m.Run()
}

type mockAuthService struct {
	pb.UnimplementedAuthServiceServer
}

// ----
// FUNC
// ----
func (s *mockAuthService) Login(_ context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email != "userEmail@gmail.com" {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}
	return &pb.LoginResponse{Token: "mock_token"}, nil
}

func (s *mockAuthService) SignUp(_ context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	if req.Username == "user" || req.Email == "userEmail@gmail.com" {
		return nil, status.Error(codes.AlreadyExists, "User exists")
	}
	return &pb.SignUpResponse{Token: "mock_token"}, nil
}

func (s *mockAuthService) ValidateJWT(_ context.Context, req *pb.JwtRequest) (*pb.JwtResponse, error) {
	if req.Token == "invalidToken" {
		return nil, status.Error(codes.InvalidArgument, "Invalid token")
	}
	return &pb.JwtResponse{
		Email: "user@example.com",
		Roles: []string{"user"},
		Exp:   1609459200,
		Iat:   1609455600,
		Nbr:   1609452000,
	}, status.Error(codes.OK, "")
}

// ----
// TESTS
// ----
func TestAuthServiceLogin(t *testing.T) {
	client := mockAuthService{}

	tests := []struct {
		name      string
		request   *pb.LoginRequest
		wantToken string
		wantErr   codes.Code
	}{
		{
			name: "Login Test #1",
			request: &pb.LoginRequest{
				Email:    "userEmail@gmail.com",
				Password: "123",
			},
			wantToken: "mock_token",
			wantErr:   codes.OK,
		},
		{
			name: "Login Test #2",
			request: &pb.LoginRequest{
				Email:    "nonUserEmail@gmail.com",
				Password: "123",
			},
			wantToken: "",
			wantErr:   codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Login(context.Background(), tt.request)
			if tt.wantErr == codes.OK {
				assert.NoError(t, err, "")
				assert.Equal(t, tt.wantToken, resp.Token)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, status.Code(err))
			}
		})
	}
}

func TestAuthServiceSignUp(t *testing.T) {
	client := mockAuthService{}

	tests := []struct {
		name      string
		request   *pb.SignUpRequest
		wantToken string
		wantErr   codes.Code
	}{
		{
			name: "SignUp Test #1",
			request: &pb.SignUpRequest{
				Username: "validUser",
				Email:    "validUserEmail@example.com",
				Password: "123",
			},
			wantToken: "mock_token",
			wantErr:   codes.OK,
		},
		{
			name: "SignUp Test #2",
			request: &pb.SignUpRequest{
				Username: "user",
				Email:    "userEmail@gmail.com",
				Password: "123",
			},
			wantToken: "",
			wantErr:   codes.AlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.SignUp(context.Background(), tt.request)
			if tt.wantErr == codes.OK {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantToken, resp.Token)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, status.Code(err))
			}
		})
	}
}

func TestAuthServiceValidate(t *testing.T) {
	cl := mockAuthService{}

	tests := []struct {
		name         string
		request      *pb.JwtRequest
		wantResponse string
		wantErr      codes.Code
	}{
		{
			name: "JwtValidation Test #1",
			request: &pb.JwtRequest{
				Token: "validToken",
			},
			wantResponse: "user@example.com",
			wantErr:      codes.OK,
		},
		{
			name: "JwtValidation Test #2",
			request: &pb.JwtRequest{
				Token: "invalidToken",
			},
			wantResponse: "",
			wantErr:      codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//resp, err := client.ValidateJWT(context.Background(), tt.request)
			resp, err := cl.ValidateJWT(context.Background(), tt.request)
			if tt.wantErr == codes.OK {
				//assert.NoError(t, err)
				assert.Equal(t, tt.wantResponse, resp.Email)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, status.Code(err))
			}
		})
	}
}
