package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/petew9527/grace"
	"net/http"
	"time"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	go func() { e.Start(":1323") }()
	grace.Wait(grace.WithOutTime(time.Second*3), grace.WithHandlers(close(e)))

}
func close(e *echo.Echo) func() error {
	return func() error {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			fmt.Println("close -> e.Shutdown ", err)
			return err
		}
		fmt.Printf("close has over\n")
		return nil
	}
}
