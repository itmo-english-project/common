package grpcclient

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/itmo-english-project/common/pkg/contexts"
	"github.com/itmo-english-project/common/pkg/grpc/grpcinterceptors"
)

type ClientConfig struct {
	target          string
	interceptors    []grpc.UnaryClientInterceptor
	dialOptions     []grpc.DialOption
	logger          *zap.Logger
	timeout         time.Duration
	enableRetry     bool
	enableTLS       bool
	keepaliveParams keepalive.ClientParameters
}

type ClientOption func(*ClientConfig)

func WithTarget(target string) ClientOption {
	return func(cfg *ClientConfig) {
		cfg.target = target
	}
}

func WithLogger(logger *zap.Logger) ClientOption {
	return func(cfg *ClientConfig) {
		cfg.logger = logger
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(cfg *ClientConfig) {
		cfg.timeout = timeout
	}
}

func WithLoggingInterceptor(ctx context.Context) ClientOption {
	return func(cfg *ClientConfig) {
		loggingOpts := []logging.Option{
			logging.WithLogOnEvents(
				logging.StartCall, logging.FinishCall,
				logging.PayloadReceived, logging.PayloadSent,
			),
		}

		interceptor := logging.UnaryClientInterceptor(
			grpcinterceptors.UnaryLoggerInterceptor(contexts.GetLogger(ctx)),
			loggingOpts...,
		)
		cfg.interceptors = append(cfg.interceptors, interceptor)
	}
}

func WithRetryInterceptor(maxRetries uint, backoff time.Duration) ClientOption {
	return func(cfg *ClientConfig) {
		retryOpts := []retry.CallOption{
			retry.WithMax(maxRetries),
			retry.WithBackoff(retry.BackoffLinear(backoff)),
		}

		interceptor := retry.UnaryClientInterceptor(retryOpts...)
		cfg.interceptors = append(cfg.interceptors, interceptor)
	}
}

func WithDefaultRetry() ClientOption {
	return WithRetryInterceptor(3, 100*time.Millisecond)
}

func WithInterceptor(interceptor grpc.UnaryClientInterceptor) ClientOption {
	return func(cfg *ClientConfig) {
		cfg.interceptors = append(cfg.interceptors, interceptor)
	}
}

func WithInsecure() ClientOption {
	return func(cfg *ClientConfig) {
		cfg.enableTLS = false
		cfg.dialOptions = append(cfg.dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
}

func WithKeepalive(time, timeout time.Duration) ClientOption {
	return func(cfg *ClientConfig) {
		cfg.keepaliveParams = keepalive.ClientParameters{
			Time:                time,
			Timeout:             timeout,
			PermitWithoutStream: true,
		}
	}
}

func WithDialOption(opt grpc.DialOption) ClientOption {
	return func(cfg *ClientConfig) {
		cfg.dialOptions = append(cfg.dialOptions, opt)
	}
}
