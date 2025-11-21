package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/health", func(c echo.Context) error {
		// Validate required connections...
		return c.String(http.StatusOK, "ok")
	})
	e.GET("/hello", handlerHello)

	terminationCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server...")
		}
	}()
	<-terminationCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}
}

func handlerHello(c echo.Context) error {
	// time.Sleep(3 * time.Second) // Test graceful shutdown
	return c.String(http.StatusOK, "Hello, SREday!")
}
