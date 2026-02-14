package grpcserver

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/itmo-english-project/common/pkg/contexts"
	"github.com/itmo-english-project/common/pkg/grpc/grpcinterceptors"
)

type Config struct {
	interceptors []grpc.UnaryServerInterceptor
	providers    []Provider
	logger       *zap.Logger
	opts         []grpc.ServerOption
}

type Option func(*Config)

func WithLogger(logger *zap.Logger) Option {
	return func(cfg *Config) {
		cfg.logger = logger
	}
}

func WithLoggingInterceptor(ctx context.Context) Option {
	return func(cfg *Config) {
		loggingOpts := []logging.Option{
			logging.WithLogOnEvents(
				logging.StartCall, logging.FinishCall,
				logging.PayloadReceived, logging.PayloadSent,
			),
		}

		interceptor := logging.UnaryServerInterceptor(
			grpcinterceptors.InterceptorLogger(contexts.GetLogger(ctx)),
			loggingOpts...)
		cfg.interceptors = append(cfg.interceptors, interceptor)
	}
}

func WithRecoveryInterceptor() Option {
	return func(cfg *Config) {
		if cfg.logger == nil {
			cfg.logger = zap.NewNop()
		}

		recoveryOpts := []recovery.Option{
			recovery.WithRecoveryHandler(func(p interface{}) (err error) {
				cfg.logger.Error("recovered from panic", zap.Any("panic", p))
				return status.Errorf(codes.Internal, "internal server error")
			}),
		}

		interceptor := recovery.UnaryServerInterceptor(recoveryOpts...)
		cfg.interceptors = append(cfg.interceptors, interceptor)
	}
}

func WithProvider(provider Provider) Option {
	return func(cfg *Config) {
		cfg.providers = append(cfg.providers, provider)
	}
}

func WithProviders(providers ...Provider) Option {
	return func(cfg *Config) {
		cfg.providers = append(cfg.providers, providers...)
	}
}

func WithInterceptor(interceptor Interceptor) Option {
	return func(cfg *Config) {
		cfg.interceptors = append(cfg.interceptors, interceptor)
	}
}

func WithServerOption(opt grpc.ServerOption) Option {
	return func(cfg *Config) {
		cfg.opts = append(cfg.opts, opt)
	}
}
