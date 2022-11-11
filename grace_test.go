package grace

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "hello\n")
	})
	go func() {
		http.ListenAndServe(":8090", nil)
	}()
	Wait(WithOutTime(time.Second*3), WithHandlers(closeHttpHandler()))
}

func closeHttpHandler() func() error {
	return func() error {
		time.Sleep(time.Second)
		fmt.Println("close success...")
		return nil
	}
}
