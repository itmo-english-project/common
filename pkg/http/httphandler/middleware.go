package httphandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"go.uber.org/zap/zapio"

	"github.com/itmo-english-project/common/pkg/contexts"
	"github.com/itmo-english-project/common/pkg/http/httpserver"
)

func DefaultMiddleware(l *zap.Logger) []httpserver.Middleware {
	var middlewares []httpserver.Middleware

	middlewares = append(middlewares,
		logger.New(logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path} ${error}\n",
			Output: &zapio.Writer{Log: l, Level: zap.DebugLevel},
		}),
		NewContext(l),
		recover.New(recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
				var errField zap.Field
				if err, ok := e.(error); ok {
					errField = zap.Error(err)
				}
				contexts.GetLogger(c.UserContext()).Error("recovered from panic in request handler", errField)
			},
		}),
	)

	return middlewares
}
