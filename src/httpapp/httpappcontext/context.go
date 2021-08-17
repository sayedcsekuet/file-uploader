package httpappcontext

import (
	"context"
	"file-uploader/src/configs"
	"file-uploader/src/crons"
	"file-uploader/src/services/fileservice"
	"file-uploader/src/services/scanner"
	"file-uploader/src/services/tokenservice"
	"fmt"
	"github.com/labstack/echo"
	logger "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RequestContext struct {
	echo.Context
	ScanInterceptor *scanner.ScanInterceptor
}

type AppContext struct {
	Server  *echo.Echo
	Setting *AppSettingContext
	scanner.ScanInterceptor
	Db *gorm.DB
	*configs.FileStorage
	fileservice.FileService
	tokenservice.TokenService
}

type AppSettingContext struct {
	Port              string
	ClamAvSocketUrl   string
	ApiKeys           []string
	DbConfig          *configs.DbConfig
	FileStorageConfig *configs.FileStorageConfig
	*crons.Cron
}

// Run starts http Server
func (a *AppContext) Run() {
	if err := a.Server.Start(fmt.Sprintf(":%v", a.Setting.Port)); err != nil && err != http.ErrServerClosed {
		a.Server.Logger.Fatal("shutting down the server")
	}
}

func (a *AppContext) WaitForInterrupt() {
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.Info("shutting down server", nil)
	if err := a.Server.Shutdown(ctx); err != nil {
		a.Server.Logger.Fatal(err)
	}
}
