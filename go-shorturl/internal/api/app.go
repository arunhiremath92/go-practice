package api

import (
	"context"
	"log"
	"net/http"
	"time"
)

type UrlShortner interface {
	ShortenUrl(string) (string, error)
	GetFullUrl(string) (string, error)
}

type AppConfig struct {
	Addr       string
	Logger     *log.Logger
	UrlService UrlShortner
}

type App struct {
	config AppConfig
	server *http.Server
}

func NewApp(appConfig AppConfig) *App {

	mux := http.NewServeMux()
	app := App{config: appConfig}
	app.registerRoutesv1(mux)

	app.server = &http.Server{
		Addr:              app.config.Addr,
		Handler:           app.withMiddleware(mux),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return &app
}

func (app *App) Start() error {
	app.config.Logger.Printf("starting http server at %s ", app.config.Addr)
	return app.server.ListenAndServe()
}

func (app *App) ShutDownServer(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
