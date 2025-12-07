package api

import (
	_ "embed" // Standard library to embed files
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// Embed the lua script into the binary so we don't need the file at runtime
//
//go:embed ratelimit.lua
var limitScript string

type RateLimiterConfig struct {
	Client *redis.Client
	Limit  int // Max burst
	Rate   int // Tokens/sec refilled
}

func RateLimitMiddleware(cfg RateLimiterConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// 1. Identify User (IP)
			// In production, use "X-Forwarded-For" if behind a proxy
			ip := r.RemoteAddr
			key := fmt.Sprintf("rate_limit:%s", ip)

			// 2. Prepare Arguments for Lua
			now := time.Now().Unix()

			// 3. Fire the Atomic Script
			// Returns 1 (Allowed) or 0 (Denied)
			result, err := cfg.Client.Eval(ctx, limitScript, []string{key},
				cfg.Limit, // ARGV[1]
				cfg.Rate,  // ARGV[2]
				now,       // ARGV[3]
				1,         // ARGV[4] (tokens requested)
			).Result()

			if err != nil {
				// Fail Open: If Redis is down, let traffic through
				// Log the error in real life
				next.ServeHTTP(w, r)
				return
			}

			// 4. Handle Result
			if result.(int64) == 0 {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("429 - Too Many Requests (Token Bucket Exhausted)\n"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
