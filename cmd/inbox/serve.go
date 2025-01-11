package main

import (
	"fmt"
	"log/slog"
	"net/http"

    "github.com/rs/cors"
)

type handlers map[string]http.Handler

func createServer(port int, origins []string, handlers handlers) *http.Server {
	mux := http.NewServeMux()
	for pattern, handler := range handlers {
		mux.Handle(pattern, handler)
	}

    c := cors.New(cors.Options{
        AllowedOrigins: origins,
        AllowedMethods: []string{http.MethodPost},
    })

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: c.Handler(mux),
	}
}

func launch(srv *http.Server, log *slog.Logger) <-chan error {
	errch := make(chan error, 1)

	go func() {
		defer close(errch)
        log.With("addr", srv.Addr).
            Info("launching server...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errch <- err
		}
	}()

	return errch
}
