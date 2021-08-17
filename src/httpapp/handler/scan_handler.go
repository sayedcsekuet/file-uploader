package handler

import (
	"file-uploader/src/helpers"
	"file-uploader/src/services/scanner"
	"fmt"
	"github.com/labstack/echo"
	logger "github.com/sirupsen/logrus"
	"net/http"
)

type ScanHandler interface {
	ScanUrls(c echo.Context) error
	ScanFiles(c echo.Context) error
	Information(c echo.Context) error
	Health(c echo.Context) error
}
type scanHandler struct {
	scanner.ScanInterceptor
	FileCollector helpers.FormCollector
}

type ScanUrlBody struct {
	Urls []string `json:"urls" validate:"required"`
}

func NewScanHandler(s scanner.ScanInterceptor, fileCollector helpers.FormCollector) ScanHandler {
	return &scanHandler{ScanInterceptor: s, FileCollector: fileCollector}
}
func (sh *scanHandler) ScanUrls(c echo.Context) error {
	body := new(ScanUrlBody)
	if err := c.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(body); err != nil {
		return err
	}
	response := sh.ScanInterceptor.ScanUrls(body.Urls)
	status := http.StatusOK
	if response.Success == false {
		status = http.StatusTeapot
	}
	return c.JSON(status, response)
}
func (sh *scanHandler) ScanFiles(c echo.Context) error {
	files, err := sh.FileCollector.CollectFromContext(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	response := sh.ScanInterceptor.ScanFiles(files)
	status := http.StatusOK
	if response.Success == false {
		status = http.StatusTeapot
	}
	return c.JSON(status, response)
}

func (sh *scanHandler) Information(c echo.Context) error {
	return c.JSON(http.StatusOK, sh.Info())
}

func (sh *scanHandler) Health(c echo.Context) error {
	if err := sh.ScanInterceptor.Scanner().Ping(); err != nil {
		logger.Error(err, nil)
		return c.JSON(
			http.StatusNotFound,
			fmt.Sprintf("Fail to connect unix socket %s", sh.ScanInterceptor.Scanner().Address()),
		)
	}
	return c.JSON(http.StatusOK, "OK")
}
