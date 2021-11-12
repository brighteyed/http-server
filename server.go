package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/brighteyed/http-server/config"
	"github.com/brighteyed/http-server/tracker"
)

func run(root *string, idleDuration *uint, port *string) {
	appCfg := config.LoadConfig("/", *root)
	if appCfg == nil {
		log.Fatal("Error loading application configuration")
	}

	for i := 0; i < len(appCfg.Locations); i++ {
		path := appCfg.Locations[i].Path
		root := appCfg.Locations[i].Root

		http.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(root))))
		log.Printf("Serving %q as %q\n", root, path)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
	}()

	idleConnsClosed := make(chan struct{})

	idleTracker := tracker.NewIdleTracker(*idleDuration)
	srv := http.Server{Addr: ":" + *port}
	srv.ConnState = idleTracker.ConnState

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
