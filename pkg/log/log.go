package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(dev bool) *zap.Logger {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	if !dev {
		encoderConfig = zap.NewProductionEncoderConfig()
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	logger := zap.New(
		zapcore.NewTee(
			zapcore.NewCore(
				encoder,
				zapcore.Lock(os.Stderr),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl >= zapcore.ErrorLevel
				}),
			),
			zapcore.NewCore(
				encoder,
				zapcore.Lock(os.Stdout),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl < zapcore.ErrorLevel && lvl >= zap.DebugLevel
				}),
			),
		),
		zap.AddCaller(),
	)

	return logger
}

func OperationField(operation string) zap.Field {
	return zap.String("operation", operation)
}
