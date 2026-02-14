package grpcclient

import (
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func NewConnection(opts ...ClientOption) (*grpc.ClientConn, error) {
	cfg := &ClientConfig{
		target:       "",
		interceptors: []grpc.UnaryClientInterceptor{},
		dialOptions:  []grpc.DialOption{},
		logger:       zap.NewNop(),
		timeout:      10 * time.Second,
		enableRetry:  true,
		enableTLS:    true,
		keepaliveParams: keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.target == "" {
		return nil, errors.New("target address is required")
	}

	cfg.dialOptions = append(cfg.dialOptions, grpc.WithKeepaliveParams(cfg.keepaliveParams))

	if len(cfg.interceptors) > 0 {
		cfg.dialOptions = append(cfg.dialOptions, grpc.WithChainUnaryInterceptor(cfg.interceptors...))
	}

	conn, err := grpc.NewClient(cfg.target, cfg.dialOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", cfg.target, err)
	}

	cfg.logger.Info("client connected", zap.String("target", cfg.target))

	return conn, nil
}

func MustConnection(logger *zap.Logger, opts ...ClientOption) *grpc.ClientConn {
	conn, err := NewConnection(opts...)
	if err != nil {
		logger.Panic("failed to create gRPC client connection", zap.Error(err))
	}
	return conn
}
