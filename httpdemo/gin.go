package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/petew9527/grace"
	"log"
	"net/http"
	"time"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Gin Hello World")
	})
	srv := &http.Server{
		Addr:    ":8090",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	grace.Wait(grace.WithOutTime(time.Second*3), grace.WithHandlers(closeGin(srv)))
}

func closeGin(srv *http.Server) func() error {
	return func() error {
		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Println("close -> gin.Shutdown ", err)
			return err
		}
		fmt.Printf("close has over\n")
		return nil
	}
}
