package errors

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"syscall"
)
const InternalServerErrorCode = "INTERNAL_SERVER_ERROR"

func CustomHTTPErrorHandler(err error, c echo.Context) {
	var code int
	var msg interface{}
	causeErr := errors.Cause(err)
	switch causeErr.(type) {
	case Known:
		e := causeErr.(Known)
		code = http.StatusBadRequest
		msg = echo.Map{"message": e.Error(), "code": e.Code(), "args": e.Args()}
	case *echo.HTTPError:
		e := causeErr.(*echo.HTTPError)
		code = e.Code
		msg = echo.Map{"message": e.Error()}
	default:
		if errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET) {
			log.Warn(err, c.Request().Context())
			return
		}
		code = http.StatusInternalServerError
		msg = echo.Map{"message": "Internal Server error", "code": InternalServerErrorCode}
		log.Error(err, c.Request().Context())
	}
	if err := c.JSON(code, msg); err != nil {
		log.Error(err, c.Request().Context())
	}
}
