package contexts

import (
	"context"

	"go.uber.org/zap"

	"github.com/itmo-english-project/common/pkg/log"
)

type valueType string

const (
	operationKey valueType = "operation_name"
	loggerKey    valueType = "logger"
)

func WithValues(parent context.Context, operation string, l *zap.Logger) context.Context {
	ctx := WithLogger(parent, l.With(log.OperationField(operation)))
	ctx = WithOperationName(ctx, operation)

	return ctx
}

func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func GetLogger(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok || l == nil {
		return zap.NewNop()
	}
	return l
}

func WithOperationName(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, operationKey, operation)
}

func GetOperationName(ctx context.Context) string {
	v, ok := ctx.Value(operationKey).(string)
	if !ok || v == "" {
		return ""
	}
	return v
}
