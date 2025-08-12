package sso

import (
	"context"
	pb "github.com/weeweeshka/sso_proto/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Sso interface {
	Register(
		ctx context.Context,
		email string,
		password string,
	) (int64, error)

	Login(
		ctx context.Context,
		email string,
		password string,
		appID int32,
	) (string, error)

	RegApp(
		ctx context.Context,
		name string,
		secret string,
	) (int32, error)
}

type serverApi struct {
	pb.UnimplementedSsoServer
	api Sso
}

func RegisterServer(grpcServer *grpc.Server, api Sso) {
	pb.RegisterSsoServer(grpcServer, &serverApi{api: api})
}

func (s *serverApi) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}

	userID, err := s.api.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{UserId: userID}, nil

}

func (s *serverApi) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "Email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required")
	}

	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "App ID is required")
	}

	token, err := s.api.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LoginResponse{Token: token}, nil
}

func (s *serverApi) RegApp(ctx context.Context, req *pb.RegappRequest) (*pb.RegappResponse, error) {

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "Name is required")
	}

	if req.GetSecret() == "" {
		return nil, status.Error(codes.InvalidArgument, "Secret is required")
	}

	appID, err := s.api.RegApp(ctx, req.GetName(), req.GetSecret())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegappResponse{AppId: appID}, nil
}
