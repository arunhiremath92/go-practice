package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
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
	RedisStore *redis.Client
	config     AppConfig
	server     *http.Server
}

func NewApp(appConfig AppConfig) *App {

	var redisOpt *redis.Options
	mux := http.NewServeMux()
	app := App{config: appConfig}
	app.registerRoutesv1(mux)
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		fmt.Println("using local redis instance")
		redisOpt = &redis.Options{
			Addr: "redis:6379",
			DB:   0,
		}
	}

	if redisUrl != "" {

		var err error
		redisOpt, err = redis.ParseURL(redisUrl)
		if err != nil {
			os.Exit(-1)
		}

	}

	app.RedisStore = redis.NewClient(redisOpt)
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
