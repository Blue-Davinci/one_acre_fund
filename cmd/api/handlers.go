package main

import (
	"context"
	"net/http"
	"os"

	"go.uber.org/zap"
)

func (app *application) generalHandler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Default hits = 0
	var hits int64 = 0
	if app.RedisDB != nil {
		val, err := app.RedisDB.Get(context.Background(), "requests_total").Int64()
		if err == nil {
			hits = val
		} else {
			app.logger.Warn("Failed to get request counter", zap.String("error", err.Error()))
		}
	}

	// Return JSON response
	err = app.writeJSON(w, http.StatusOK, envelope{
		"hits":     hits,
		"hostName": hostname,
		"success":  true,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Log request
	app.logger.Info("Handled general request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Int64("hits", hits))
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"status": "healthy"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// logger
	app.logger.Info("Handled health check request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path))
}
