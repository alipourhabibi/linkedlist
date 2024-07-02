package main

import (
	"context"
	"errors"
	"flag"
	"linkedlist/api"
	"linkedlist/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var cfg = flag.String("config", "config/config.yaml", "yaml config path")

func main() {

	flag.Parse()

	err := config.Load(*cfg)
	if err != nil {
		slog.Error("Reading configuration", "error", err)
		os.Exit(1)
	}
	level, ok := config.MapLevel[strings.ToUpper(config.Confs.Logger.Level)]
	if !ok {
		level = slog.LevelError
	}
	l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: config.Confs.Logger.AddSource,
		Level:     level,
	}))
	slog.SetDefault(l)

	server, err := api.New()
	if err != nil {
		slog.Error("Bootin api", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	go func() {
		err = server.Start(ctx)
		if err != nil {
			return
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	for {
		sig := <-sigs
		switch sig {
		case syscall.SIGHUP:
			slog.Info("Received SIGHUP, reloading configuration...")

			// reload the config
			err := config.Load("config/config.yaml")
			if err != nil {
				slog.Error("Reading configuration", "error", err)
				os.Exit(1)
				return
			}

			// reload the logger with new config
			level, ok := config.MapLevel[strings.ToUpper(config.Confs.Logger.Level)]
			if !ok {
				level = slog.LevelError
			}
			l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				AddSource: config.Confs.Logger.AddSource,
				Level:     level,
			}))
			slog.SetDefault(l)

			// shutdown the prev server
			// continue with the the last config if failed
			err = server.Shutdown(ctx)
			if err != nil {
				slog.Error("could not stop server", "error", err)
				continue
			}

			server, err = api.New()
			if err != nil {
				slog.Error("Bootin api", "error", err)
				os.Exit(1)
			}

			go func() {
				// TODO read about this error
				if err := server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
					slog.Error("could not start new server", "error", err)
					os.Exit(1)
				}
			}()

		case syscall.SIGINT, syscall.SIGTERM:
			slog.Info("Received SIGINT/SIGTERM, shutting down...")
			if err := server.Shutdown(ctx); err != nil {
				slog.Error("could not gracefully shut down server", "error", err)
			}
			os.Exit(0)
		}
	}

}
