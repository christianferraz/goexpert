package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	server := &http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(4 * time.Second)
		w.Write([]byte("Hello\n"))
	})
	go func() {
		fmt.Printf("Server is running at %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	fmt.Println("Server Exited Properly")
}