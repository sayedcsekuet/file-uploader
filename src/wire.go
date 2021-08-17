// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"file-uploader/src/httpapp"
	"file-uploader/src/httpapp/httpappcontext"
	"github.com/google/wire"
)

func InitializeApp(
	appSetting httpappcontext.AppSettingContext,
) (*httpappcontext.AppContext, error) {
	wire.Build(
		httpapp.Initialize,
	)

	return &httpappcontext.AppContext{}, nil
}
