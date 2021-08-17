package router

import (
	"file-uploader/src/helpers"
	"file-uploader/src/httpapp/handler"
	"file-uploader/src/httpapp/httpappcontext"
	"file-uploader/src/httpapp/middleware"
)

type Router struct {
	AppContext *httpappcontext.AppContext
}

func NewRouter(appContext *httpappcontext.AppContext) *Router {
	return &Router{AppContext: appContext}
}
func (r *Router) Init() {
	r.InitScanner()
	r.InitFiles()
}

func (r *Router) InitScanner() {
	reqHandler := handler.NewScanHandler(r.AppContext.ScanInterceptor, helpers.NewFileCollector())
	r.AppContext.Server.GET("/health", reqHandler.Health)
	g := r.AppContext.Server.Group("/scanners", middleware.ApiKeyMiddleware(r.AppContext.Setting.ApiKeys))
	g.GET("/info", reqHandler.Information)
	g.POST("/files", reqHandler.ScanFiles)
	g.POST("/urls", reqHandler.ScanUrls)
}

func (r *Router) InitFiles() {
	// Initialize handlers
	reqHandler := handler.NewFileHandler(
		helpers.NewFileCollector(),
		r.AppContext.FileService,
	)
	tokenHandler := handler.NewTokenHandler(r.AppContext.TokenService)
	// Initialize routes
	g := r.AppContext.Server.Group("/files", middleware.ApiKeyMiddleware(r.AppContext.Setting.ApiKeys))
	g.GET("", reqHandler.List)
	g.POST("", reqHandler.Create)
	g.GET("/:id", reqHandler.Get)
	g.GET("/:id/download", reqHandler.Download)
	g.GET("/:id/stream", reqHandler.Stream)
	g.DELETE("/:id", reqHandler.Delete)

	g.POST("/tokens", tokenHandler.GenerateTokens)
}
