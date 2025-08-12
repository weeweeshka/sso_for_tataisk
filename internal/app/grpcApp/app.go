package grpcApp

import (
	"fmt"
	grpcHandlers "github.com/weeweeshka/sso_for_tataisk/internal/grpc/sso"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	gRPCServer *grpc.Server
	logr       *zap.Logger
	port       int
}

func New(port int, logr *zap.Logger, sso grpcHandlers.Sso) *GRPCServer {

	gRPCServer := grpc.NewServer()
	grpcHandlers.RegisterServer(gRPCServer, sso)

	return &GRPCServer{gRPCServer, logr, port}
}

func (s *GRPCServer) Run() error {

	l, err := net.Listen("tcp", ":"+string(s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	if err := s.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}

func (s *GRPCServer) GracefulStop() error {
	s.gRPCServer.GracefulStop()
	return nil
}
