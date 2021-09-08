package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ilyazamyslov/inet-scanner-golang/internal/handler"
	"github.com/ilyazamyslov/inet-scanner-golang/internal/repository"
	"github.com/ilyazamyslov/inet-scanner-golang/internal/service"
	riaken_core "github.com/riaken/riaken-core"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)

	r := chi.NewRouter()
	addrs := []string{"127.0.0.1:8087"}
	client := riaken_core.NewClient(addrs, 1)
	if err := client.Dial(); err != nil {
		logger.Fatal().Err(err).Msg("DB initializing error") //(err.Error()) // all nodes are down
	}
	defer client.Close()
	session := client.Session()
	defer session.Release()

	dbRepo := &repository.DB{DB: session}
	service := service.NewScannerService(&logger, dbRepo)
	//h := handler.New(&logger, service)
	scanHostHandler := handler.NewHostScan(&logger, service)
	scanNetworkHandler := handler.NewNetworkScan(&logger, service)

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.RequestLogger(&handler.LogFormatter{Logger: &logger}))
		r.Use(middleware.Recoverer)
		r.Method(http.MethodGet, handler.ScanHostPath, scanHostHandler)
		r.Method(http.MethodGet, handler.ScanNetworkPath, scanNetworkHandler)
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		logger.Info().Msgf("Server is listening on :%d", 8080)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server error")
		}
	}()
	<-shutdown

	logger.Info().Msg("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server shutdown error")
	}

	logger.Info().Msg("Server stopped gracefully")
}
