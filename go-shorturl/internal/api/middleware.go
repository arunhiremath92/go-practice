package api

import (
	"log"
	"net/http"
	"time"
)

func (app *App) withMiddleware(next http.Handler) http.Handler {
	limiter := RateLimitMiddleware(RateLimiterConfig{
		Client: app.RedisStore,
		Limit:  10,
		Rate:   1,
	})

	return limiter(app.recoveryMiddleware(app.loggingMiddleware(next)))
}

func (a *App) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		a.config.Logger.Printf("started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		a.config.Logger.Printf("completed %s %s in %s",
			r.Method, r.URL.Path, time.Since(start))
	})
}

func (a *App) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %+v", rec)
				writeError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
