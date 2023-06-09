package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

func newRouter() *httprouter.Router {
	mux := httprouter.New()
	ytApiKey := os.Getenv("YOUTUBE_API_KEY")
	ytChannelId := os.Getenv("YOUTUBE_CHANNEL_ID")

	if ytApiKey == "" {
		log.Fatal("youtube API key is not provided")
	}

	if ytApiKey == "" {
		log.Fatal("youtube ChannelId is not provided")
	}

	mux.GET("/ytoutube/channel/stats", getChannelStats(ytApiKey, ytChannelId))
	return mux
}

func main() {
	srv := &http.Server{
		Addr:    ":11000",
		Handler: newRouter(),
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		log.Println("service interrupt received")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown error %v", err)
		}

		log.Println("shutdown complete")

		close(idleConnsClosed)
	}()

	log.Println("Starting server on port 11000")

	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("fatal http server failed to start: %v", err)
		}
	}

	<-idleConnsClosed
	log.Println("Service Stop")
}
