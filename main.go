package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/brighteyed/http-server/config"
	"github.com/brighteyed/http-server/tracker"
)

func main() {
	root := flag.String("d", ".", "the directory of files to host")
	port := flag.String("p", "8100", "port to serve on")
	idleDuration := flag.Uint("t", 0, "duration before shutdown while inactive (0 – disable)")
	flag.Parse()

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

	idleConnsClosed := make(chan struct{})

	idleTracker := tracker.NewIdleTracker(*idleDuration)
	srv := http.Server{Addr: ":" + *port}
	srv.ConnState = idleTracker.ConnState

	go func() {
		<-idleTracker.Done()

		log.Println("Shutting down")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Shutdown, %v", err)
		}

		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServer, %v", err)
	}

	// Wait until Shutdown returns
	<-idleConnsClosed
}