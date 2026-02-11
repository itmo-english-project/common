package httphandler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/itmo-english-project/exchange/pkg/contexts"
	"github.com/itmo-english-project/exchange/pkg/http/httpserver"
)

func NewContext(l *zap.Logger) httpserver.Middleware {
	return func(c *fiber.Ctx) error {
		logger := l
		var ctx context.Context = c.Context()
		ctx = contexts.WithValues(
			ctx,
			c.Path(),
			logger,
		)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
