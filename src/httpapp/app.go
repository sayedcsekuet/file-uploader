package httpapp

import (
	"file-uploader/src/configs"
	"file-uploader/src/crons"
	appErrors "file-uploader/src/errors"
	"file-uploader/src/helpers"
	"file-uploader/src/httpapp/httpappcontext"
	"file-uploader/src/httpapp/router"
	"file-uploader/src/repositories"
	"file-uploader/src/services/fileservice"
	"file-uploader/src/services/filestorage"
	"file-uploader/src/services/scanner"
	"file-uploader/src/services/tokenservice"
	"file-uploader/src/utils"
	"file-uploader/src/validators"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/newrelic/go-agent/_integrations/nrecho"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"runtime"
)

// Initialize instantiates http application with configured routes
func Initialize(appSetting httpappcontext.AppSettingContext, ) *httpappcontext.AppContext {
	e := echo.New()
	e.Use(nrecho.Middleware(utils.NewRelicApp))
	context := InitializeAppContext(&appSetting, e)
	// Registering validator
	e.Validator = validators.NewAppValidator()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"method":"${method}", "uri":"${uri}", "message":"Receive request",` +
			`"context":{"correlation_id":"${header:X-Correlation-Id}","tag":"file-uploader",` +
			`"user_agent":"${user_agent}","error":"${error}","latency":${latency},` +
			`"latency_human":"${latency_human}"},"level_name":"INFO","level":${status},` +
			`"datetime_iso":"${time_rfc3339_nano}"}` + "\n",
	}))
	e.HTTPErrorHandler = appErrors.CustomHTTPErrorHandler
	//Initialize router
	appRouter := router.NewRouter(context)
	appRouter.Init()

	return context
}

func InitializeAppContext(appSetting *httpappcontext.AppSettingContext, e *echo.Echo) *httpappcontext.AppContext {

	db, err := configs.InitDB(appSetting.DbConfig)
	if err != nil {
		err = errors.Wrap(err, "Could not connect to Mysql")
		logger.Fatal(err, nil)
	}

	fileStorage := configs.NewFileStorage(appSetting.FileStorageConfig)
	interceptor := scanner.NewScanInterceptor(new(scanner.Clamav), appSetting.ClamAvSocketUrl)
	tokenService := tokenservice.NewTokenService()
	fileService := fileservice.NewFileService(
		interceptor,
		repositories.NewFileRepository(db),
		filestorage.NewStorageService(fileStorage),
		tokenService,
	)
	context := &httpappcontext.AppContext{
		Server:          e,
		Setting:         appSetting,
		ScanInterceptor: interceptor,
		FileStorage:     fileStorage,
		Db:              db,
		FileService:     fileService,
		TokenService:    tokenService,
	}
	// Add cron job
	addCronJobs(appSetting.Cron, fileService)
	return context
}

func addCronJobs(c *crons.Cron, fileService fileservice.FileService) {
	// Add auto deleting expired files jobs
	c.Add(func() (*gocron.Job, error) {
		cronStr := fmt.Sprintf("0/%d * * * *", helpers.RandomInt(10, 20))
		logger.Infof("Cron job starting pattern: %s", cronStr)
		j, err := c.Cron(cronStr).Do(func() {
			if runtime.GOOS == "windows" {
				return
			}
			//Refresh the calmav databse
			err := fileService.DeleteExpiredFiles()
			if err != nil {
				logger.Infof("Deleting expired files Error: %v", err)
			}
		})
		return j, err
	})
}
