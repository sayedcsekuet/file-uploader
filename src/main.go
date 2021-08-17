package main

import (
	"file-uploader/src/configs"
	"file-uploader/src/crons"
	"file-uploader/src/helpers"
	"file-uploader/src/httpapp/httpappcontext"
	"file-uploader/src/utils"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func init() {
	// If FILE_UPLOADER_API_KEY is not set try to read from file and set it to the env
	if os.Getenv("FILE_UPLOADER_API_KEY") == "" {
		b, _ := ioutil.ReadFile(os.Getenv("FILE_UPLOADER_API_KEY_FILE"))
		_ = os.Setenv("FILE_UPLOADER_API_KEY", string(b))
	}
	rEnv := []string{"ENV_NAME", "FILE_UPLOADER_API_KEY", "DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "FILE_STORAGE_DRIVER", "FILE_STORAGE_REGION"}
	for _, e := range rEnv {
		if os.Getenv(e) == "" {
			m := errors.Wrap(errors.New(e), "Environment value will not be empty!")
			log.Fatal(m, nil)
		}
	}
	// Init logger
	InitLogger()
}

func InitLogger() {
	formatter := &log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyMsg:   "message",
			log.FieldKeyLevel: "level_name",
			log.FieldKeyTime:  "datetime_iso",
		},
	}
	log.SetFormatter(formatter)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
func main() {
	go func() {
		if err := utils.InitNewRelicApp(); err != nil {
			err = errors.Wrap(err, "[✗] Couldn't instantiate NewRelic")
			log.Error(err, nil)
		}
	}()
	httpPort := os.Getenv("HTTP_PORT")

	clamUrl := os.Getenv("CLAMAV_SOCKET_URL")
	if !helpers.IsValidUrl(clamUrl) {
		log.Fatal("ClamAv env variable 'CLAMAV_SOCKET_URL' value is not valid url", nil)
	}
	dbConfig := configs.DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DbName:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Pass:     os.Getenv("DB_PASS"),
		PassFile: os.Getenv("DB_PASS_FILE"),
	}
	fileStorageConfig := configs.FileStorageConfig{
		Provider:         os.Getenv("FILE_STORAGE_DRIVER"),
		Endpoint:         os.Getenv("FILE_STORAGE_ENDPOINT"),
		DisableSSL:       os.Getenv("FILE_STORAGE_DISABLE_SSL") != "",
		Region:           os.Getenv("FILE_STORAGE_REGION"),
		AccessKeyID:      os.Getenv("FILE_STORAGE_ACCESS_KEY"),
		SecretAccessKey:  os.Getenv("FILE_STORAGE_SECRET"),
		Profile:          os.Getenv("FILE_STORAGE_PROFILE"),
		S3ForcePathStyle: os.Getenv("S3_FORCE_PATH_STYLE") == "",
	}
	cron := crons.NewCron()

	appSetting := httpappcontext.AppSettingContext{
		Port:              httpPort,
		ClamAvSocketUrl:   clamUrl,
		ApiKeys:           helpers.CollectApiKeys(),
		DbConfig:          &dbConfig,
		FileStorageConfig: &fileStorageConfig,
		Cron:              cron,
	}
	appContext, err := InitializeApp(appSetting)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Couldn't initialise HTTP app"), nil)
	}

	go appContext.Run()

	// Start cron
	cron.StartScheduler()

	appContext.WaitForInterrupt()

	// Gracefully stop important processes if any

	log.Info("exited", nil)
}
