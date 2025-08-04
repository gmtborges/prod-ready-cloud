package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/hello", handlerHello)

	e.Logger.Fatal(e.Start(":1323"))
}

func handlerHello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, SREday!")
}
