package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NouNhaN-GitHub/assessment-tax/ktaxes"
	"github.com/NouNhaN-GitHub/assessment-tax/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	p, err := postgres.New()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	handler := ktaxes.New(p)

	ktax := e.Group("/ktaxes")
	{
		ktax.GET("/allowance", handler.AllowanceHandler)
	}

	tax := e.Group("/tax")
	{
		tax.POST("/calculations", handler.TaxCalculationsHandler)
	}

	admin := e.Group("/admin")
	admin.Use(middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		if username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD") {
			return true, nil
		}

		return false, nil
	}))
	{
		admin.POST("/deductions/personal", handler.PersonalDeductionHandler)
	}

	go func() {
		if err := e.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	fmt.Println("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
