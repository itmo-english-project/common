package grpcserver

import (
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Provider interface {
	RegisterGRPCServer(s *grpc.Server)
}

type Interceptor = grpc.UnaryServerInterceptor

func New(opts ...Option) *grpc.Server {
	cfg := &Config{
		interceptors: []grpc.UnaryServerInterceptor{},
		providers:    []Provider{},
		logger:       zap.NewNop(),
		opts:         []grpc.ServerOption{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if len(cfg.interceptors) > 0 {
		cfg.opts = append(cfg.opts, grpc.ChainUnaryInterceptor(cfg.interceptors...))
	}

	server := grpc.NewServer(cfg.opts...)

	for _, provider := range cfg.providers {
		provider.RegisterGRPCServer(server)
	}

	return server
}

func MustListenAndServe(server *grpc.Server, port int, logger *zap.Logger) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Panic("failed to listen", zap.Error(err), zap.Int("port", port))
	}

	logger.Info("gRPC server started", zap.Int("port", port))
	if err = server.Serve(lis); err != nil {
		logger.Panic("failed to serve", zap.Error(err))
	}
}

func ListenAndServe(server *grpc.Server, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %w", port, err)
	}

	return server.Serve(lis)
}
