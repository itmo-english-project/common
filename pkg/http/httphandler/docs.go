package httphandler

import (
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

type DocsProvider struct {
	fs fs.FS
}

func NewDocsProvider(fs fs.FS) *DocsProvider {
	return &DocsProvider{fs: fs}
}

func (p *DocsProvider) RegisterFastHTTPRoutes(a fiber.Router) {
	a.Use("/docs", filesystem.New(filesystem.Config{
		Root:   http.FS(p.fs),
		Browse: true,
	}))
}
