package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"rina.icu/pocman-server/internal/conf"
	"rina.icu/pocman-server/internal/handler"
	"rina.icu/pocman-server/internal/server"
)

var config_file = flag.String("c", "config.yml", "Configuration file, example: config.yml")

func main() {
	// parse cmd parameters
	flag.Parse()

	// init config
	log.Printf("Reading configuration from \"%s\"", *config_file)
	conf.AppConfig.Load(*config_file)

	// init HTTP server
	if conf.AppConfig.Server.Mode == "release" {
		// set gin log level
		gin.SetMode(gin.ReleaseMode)
	}

	pingHandler := handler.NewPingHandler()
	r := server.New(pingHandler)

	addr := fmt.Sprintf("%s:%d", conf.AppConfig.Server.Host, conf.AppConfig.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// run HTTP server
	go func() {
		log.Printf("Starting Pocman server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// block and wait
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Pocman server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Pocman server forced to shutdown:", err)
	}

	log.Println("Pocman server exited")
}
