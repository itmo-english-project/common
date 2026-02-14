package grpcinterceptors

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
)

func UnaryLoggerInterceptor(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			if keyStr, ok := key.(string); ok {
				f = append(f, zap.Any(keyStr, value))
			}
		}

		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg, f...)
		case logging.LevelInfo:
			l.Info(msg, f...)
		case logging.LevelWarn:
			l.Warn(msg, f...)
		case logging.LevelError:
			l.Error(msg, f...)
		}
	})
}
