package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHello(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	expectedBody := "Hello, SREday!"

	err := handlerHello(c)
	if err != nil {
		t.Errorf("Error %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Status code: %d, want 200", rec.Code)
	}
	if rec.Body.String() != expectedBody {
		t.Errorf("Exptected 'Hello, SREDay!', but got %s", rec.Body.String())
	}
}
