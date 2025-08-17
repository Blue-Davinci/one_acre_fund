package main

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or
			// not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the
				// response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been
				// sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has the type any, so we use
				// fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. In turn, this will log the error using
				// our custom Logger type at the ERROR level and send the client a 500
				// Internal Server Error response.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// incrementorMiddleware increments request count in Redis for each request.
func (app *application) incrementorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		if app.RedisDB != nil {
			err := app.RedisDB.Incr(ctx, "requests_total").Err()
			if err != nil {
				app.logger.Warn("Failed to increment request counter", zap.String("error", err.Error()))
			} else {
				val, _ := app.RedisDB.Get(ctx, "requests_total").Result()
				w.Header().Set("X-Request-Count", val) // optional: add to response
			}
		}
		next.ServeHTTP(w, r)
	})
}
