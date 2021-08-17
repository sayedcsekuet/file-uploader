package handler

import (
	"file-uploader/src/services/tokenservice"
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/labstack/echo"
	logger "github.com/sirupsen/logrus"
	"net/http"
)

type TokenHandler interface {
	GenerateTokens(c echo.Context) error
}
type tokenHandler struct {
	tokenservice.TokenService
}

func NewTokenHandler(ts tokenservice.TokenService) TokenHandler {
	return &tokenHandler{ts}
}

type TokenBody struct {
	Ids       []string `json:"ids" validate:"required"`
	ExpiredAt string   `json:"expired_at"`
}

func (jh *tokenHandler) GenerateTokens(c echo.Context) error {
	ownerId := c.Request().Header.Get("x-api-key")
	body := new(TokenBody)
	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var response = map[string]string{}
	var exp int64 = 0
	if body.ExpiredAt != "" {
		t, err := dateparse.ParseAny(body.ExpiredAt)
		if err != nil {
			logger.Error(err, nil)
			return echo.NewHTTPError(http.StatusBadRequest, "Expired at time is not valid!")
		}
		exp = t.Unix()
	}
	for _, id := range body.Ids {
		token, err := jh.TokenService.Generate(id, ownerId, exp)
		if err != nil {
			logger.Error(err, nil)
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Fail to generate token with id %s!", id))
		}
		response[id] = token
	}
	return c.JSON(http.StatusCreated, response)
}
