package internalgrpc

import (
	"fmt"
	"net"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	pb "github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/pb/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer *grpc.Server
	logger     logger.Logger
}

func NewServer(logger logger.Logger, eventHandler pb.EventsServer) *Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(LoggingInterceptor),
	)
	pb.RegisterEventsServer(s, eventHandler)

	reflection.Register(s)

	return &Server{logger: logger, grpcServer: s}
}

func (s *Server) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.logger.Info("GRPC server start on " + addr)

	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	s.grpcServer.GracefulStop()
	s.logger.Info("GRPC server stopped")
}
