package httphandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type SwaggerProvider struct {
	path string
}

func NewSwaggerProvider(path string) *SwaggerProvider {
	return &SwaggerProvider{path: path}
}

func (p *SwaggerProvider) RegisterFastHTTPRoutes(r fiber.Router) {
	r.Get("/swagger/*", swagger.New(swagger.Config{
		URL: p.path,
	}))
}
