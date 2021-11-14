package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/brighteyed/http-server/config"
	"github.com/brighteyed/http-server/server"
	"github.com/brighteyed/http-server/tracker"
)

func main() {
	root := flag.String("d", "", "the directory of files to host")
	port := flag.String("p", "8100", "port to serve on")
	idle := flag.Uint("t", 0, "duration before shutdown while inactive (0 – disable)")
	flag.Parse()

	run(*root, *port, *idle)
}

func run(root string, port string, idleDuration uint) {
	appCfg := config.NewAppConfig(root)
	if appCfg == nil {
		log.Fatal("Error loading application configuration")
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
	}()

	idleTracker := tracker.NewIdleTracker(idleDuration)
	idleConnsClosed := make(chan struct{})

	srv := http.Server{
		Addr:      ":" + port,
		Handler:   server.NewHandler(appCfg.Locations),
		ConnState: idleTracker.ConnState,
	}

	go func() {
		select {
		case <-idleTracker.Done():
			shutdown(&srv)
		case <-signalChan:
			shutdown(&srv)
		}

		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServer, %v", err)
	}

	// Wait until Shutdown returns
	<-idleConnsClosed
}

func shutdown(srv *http.Server) {
	log.Println("Shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Shutdown, %v", err)
	}
}
