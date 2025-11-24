package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("No .env file found, using environment variables")
	}
	logLevel := os.Getenv("LOG_LEVEL")
	switch strings.ToUpper(logLevel) {
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	terminationCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	shutdownMeterProvider, err := initMeterProvider(terminationCtx)
	shutdownTraceProvider, err := initTracerProvider(terminationCtx)

	e := echo.New()
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	e.Use(otelecho.Middleware(serviceName))

	e.GET("/health", func(c echo.Context) error {
		// Validate required connections...
		return c.String(http.StatusOK, "ok")
	})
	e.GET("/hello", handlerHello)

	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("shutting down the server...")
		}
	}()
	<-terminationCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = errors.Join(
		e.Shutdown(shutdownCtx),
		shutdownMeterProvider(shutdownCtx),
		shutdownTraceProvider(shutdownCtx),
	)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msg("Server shutdown successfully.")
}

func handlerHello(c echo.Context) error {
	// time.Sleep(3 * time.Second) // Test graceful shutdown
	return c.String(http.StatusOK, "Hello, SREday!")
}
