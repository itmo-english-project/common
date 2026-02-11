package httpserver

import "github.com/gofiber/fiber/v2"

type Provider interface {
	RegisterFastHTTPRoutes(r fiber.Router)
}

type Middleware = fiber.Handler

func New(middlewares []Middleware, providers []Provider, errHandler fiber.ErrorHandler) *fiber.App {
	config := fiber.Config{
		DisableStartupMessage: true,
		ReadBufferSize:        16 * 1024,
		ErrorHandler:          errHandler,
		Immutable:             true,
	}

	server := fiber.New(config)
	for _, m := range middlewares {
		server.Use(m)
	}

	for _, p := range providers {
		p.RegisterFastHTTPRoutes(server)
	}

	return server
}
