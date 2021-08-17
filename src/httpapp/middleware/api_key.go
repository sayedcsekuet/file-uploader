package middleware

import (
	"file-uploader/src/helpers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func ApiKeyMiddleware(keys []string) echo.MiddlewareFunc {
	// Validating api key using middleware
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:x-api-key",
		Validator: func(key string, c echo.Context) (bool, error) {
			return helpers.Contains(keys, key), nil
		},
		Skipper: func(context echo.Context) bool {
			if context.QueryParam("token") != "" {
				return true
			}
			return false
		},
	})
}
