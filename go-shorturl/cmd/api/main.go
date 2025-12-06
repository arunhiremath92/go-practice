package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/arunhiremath92/go-shorturl/internal/api"
	"github.com/arunhiremath92/go-shorturl/internal/urlshortner"
	store "github.com/arunhiremath92/go-shorturl/pkg/redis"
	"github.com/redis/go-redis/v9"
)

func main() {
	loghandler := log.New(os.Stdout, "url-shortner:", log.Default().Flags())
	addr := ":8000"
	var redisOpt *redis.Options
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		loghandler.Println("using local redis instance")
		redisOpt = &redis.Options{
			Addr: "redis:6379",
			DB:   0,
		}
	}

	if redisUrl != "" {
		loghandler.Println("using remote redis")
		var err error
		redisOpt, err = redis.ParseURL(redisUrl)
		if err != nil {
			loghandler.Printf("failed to set up redis connection %s", err)
			os.Exit(-1)
		}

	}

	redisStore := store.NewRedisStore(redisOpt)
	urlShortenerSvc := urlshortner.NewUrlShortner(redisStore)

	app := api.NewApp(api.AppConfig{Addr: addr, Logger: loghandler,
		UrlService: urlShortenerSvc})

	// handle errors from the http
	appErrorHandler := make(chan error, 1)
	// setup interruption Handler
	signalHandler := make(chan os.Signal, 1)
	signal.Notify(signalHandler, os.Interrupt)

	go func() {
		appErrorHandler <- app.Start()
	}()

	select {
	case serverError := <-appErrorHandler:
		loghandler.Printf("server error: %v", serverError)
	case <-signalHandler:
		loghandler.Print("received interruption from the os")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := app.ShutDownServer(ctx); err != nil {
			loghandler.Printf("graceful shutdown failed: %v", err)
		}
	}

}
