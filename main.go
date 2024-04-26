package main

import (
	"net/http"

	"github.com/NouNhaN-GitHub/assessment-tax/ktaxes"
	"github.com/NouNhaN-GitHub/assessment-tax/postgres"
	"github.com/labstack/echo/v4"
)

func main() {
	p, err := postgres.New()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	handler := ktaxes.New(p)
	e.GET("/", handler.AllowanceHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
