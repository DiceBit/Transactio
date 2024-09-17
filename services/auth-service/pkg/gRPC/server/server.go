package server

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"log"
	"net"
	"user-service/pkg/db/pgx"
	userUtils "user-service/pkg/db/utils"
	pb "user-service/pkg/gRPC/proto"
	"user-service/pkg/utils"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
	db *pgxpool.Pool
}

func (authServ *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	usr, err := userUtils.UsrByEmail(ctx, authServ.db, req.Email)
	if err != nil {
		log.Println("Error with getting usr by email")
		return nil, err
	}

	if ok := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(req.Password)); ok != nil {
		log.Println("Invalid password")
		return nil, nil
	}

	token, err := utils.GenerateJWT(usr.Email, usr.Role)
	if err != nil {
		log.Println("Error generating token:", err)
		return nil, err
	}

	return &pb.LoginResponse{Token: token}, nil
}

// registration
func (authServ *AuthServiceServer) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	if exist, err := userUtils.CheckIfExistUsr(ctx, authServ.db, req); !exist && err == nil {
		err = userUtils.AddUser(ctx, authServ.db, req)
		if err != nil {
			log.Println("Error with adding usr.", err)
		}
		log.Println("User added")
	} else if exist && err == nil {
		log.Println("User already exists")
		return nil, nil
	} else {
		log.Println("Error with checking usr in DB.", err)
	}

	usr, err := userUtils.UsrByEmail(ctx, authServ.db, req.Email)
	if err != nil {
		log.Println("Error with getting usr by email")
		return nil, err
	}
	token, err := utils.GenerateJWT(usr.Email, usr.Role)
	if err != nil {
		log.Println("Error generating token:", err)
		return nil, err
	}

	return &pb.SignUpResponse{Token: token}, nil
}

func (authServ *AuthServiceServer) ValidateJWT(ctx context.Context, req *pb.JwtRequest) (*pb.JwtResponse, error) {
	claims, err := utils.ValidateJWT(req.Token)
	if err != nil {
		log.Println("Error with validation jwt.", err)
		return nil, err
	}

	exp := claims.ExpiresAt.Unix()
	iat := claims.IssuedAt.Unix()
	nbr := claims.NotBefore.Unix()

	return &pb.JwtResponse{Email: claims.Email, Roles: claims.Roles,
		Exp: exp, Iat: iat, Nbr: nbr}, nil
	/*return &pb.JwtResponse{Email: claims.Email, Roles: claims.Roles,
	Exp: exp}, nil*/
}

func RunServe(authName, authServiceAddr string) {
	listen, err := net.Listen("tcp", authServiceAddr)
	if err != nil {
		log.Fatalf("Error with starting %s. %v", authName, err)
	}

	grpcServer := grpc.NewServer()

	log.Printf("%s is running on %s", authName, authServiceAddr)

	newAuthServ := NewAuthService()

	pb.RegisterAuthServiceServer(grpcServer, newAuthServ)
	log.Fatal(grpcServer.Serve(listen))
}

func NewAuthService() *AuthServiceServer {
	db, err := pgx.New()
	if err != nil {
		log.Println("Error with DB", err)
	}

	authServ := AuthServiceServer{
		db: db,
	}
	return &authServ
}
