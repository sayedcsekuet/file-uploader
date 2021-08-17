package utils

import (
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrlogrus"
	"os"
)

var NewRelicApp newrelic.Application

func InitNewRelicApp() error {
	license := os.Getenv("NEWRELIC_APM_LICENSE")
	if license == "" {
		return nil
	}
	appName := os.Getenv("NEWRELIC_APPNAME_PREFIX") + "-" + os.Getenv("NEWRELIC_HOST_DISPLAY_NAME")
	config := newrelic.NewConfig(appName, license)
	config.Logger = nrlogrus.StandardLogger()
	app, err := newrelic.NewApplication(config)
	NewRelicApp = app

	return err
}
