package httphandler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/itmo-english-project/exchange/pkg/contexts"
	"github.com/itmo-english-project/exchange/pkg/http/httperrors"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	ctx := c.UserContext()
	code := c.Response().StatusCode()

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.JSON(httperrors.Error{
		Message: err.Error(),
		Code:    strconv.Itoa(code),
		Details: map[string]interface{}{
			"operation_name": contexts.GetOperationName(ctx),
		},
	})
}
