package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	dumpLog = log.New(os.Stdout, "", 0)
	errLog  = log.New(os.Stderr, "[ERROR]: ", log.LstdFlags)
	infoLog = log.New(os.Stdout, "[INFO]: ", log.LstdFlags)
)

var (
	addr string
	body bool
)

func main() {
	flag.StringVar(&addr, "addr", ":8080", "addr optionally specifies the TCP address for the server to listen on, in the form \"host:port\"")
	flag.BoolVar(&body, "body", true, "If true, also dump the request body")
	flag.Parse()

	s := http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dump, err := httputil.DumpRequest(r, body)
			if err != nil {
				errLog.Printf("Unable to dump request: %s", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// Loggers are OK for concurrent use from multiple goroutines.
			dumpLog.Printf("%s\n", dump)
		}),
	}

	errCh := make(chan error, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		infoLog.Printf("Accepting connections on %q\n\n", addr)
		errCh <- s.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		errLog.Printf("Exiting with error: %s", err)
	case sig := <-sigCh:
		infoLog.Printf("Exiting with signal: %s. Waiting for in-flight connections to finish...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			errLog.Printf("Unable to shutdown server: %s", err)
		}
	}
}
