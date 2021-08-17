package httpapp

import (
	"file-uploader/src/errors"
	"github.com/labstack/echo"
	"net/http"
)

func ErrorResponse(err error) *echo.HTTPError {
	switch err.(type) {
	case errors.Known:
		e := err.(errors.Known)
		return echo.NewHTTPError(e.Code(), err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
}
