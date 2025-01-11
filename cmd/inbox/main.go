package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/bdreece/inbox/pkg/controller"
	"github.com/bdreece/inbox/pkg/email"
	"github.com/bdreece/inbox/pkg/logger"
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New(
		validator.WithRequiredStructEnabled(),
	)
)

func main() {
	defer quit()

	port := flag.Int("p", 3000, "port")
	env := flag.String("e", "dev", "environment")
	configs := flag.String("c", "configs", "config dir")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := loadConfig(
		filepath.Join(*configs, "settings.yml"),
		filepath.Join(*configs, fmt.Sprintf("settings.%s.yml", *env)),
	)
	if err != nil {
		panic(err)
	}

	logOpts, err := logger.Configure(cfg)
	if err != nil {
		panic(err)
	}
	log, err := logger.New(logOpts)
	if err != nil {
		panic(err)
	}

	sesOpts, err := email.ConfigureSES(cfg)
	if err != nil {
		panic(err)
	}
	client, err := email.NewSESClient(ctx, *sesOpts)
	if err != nil {
		panic(err)
	}

    corsOpts := new(struct{
        Origins []string `yaml:"origins"`
    })
    if err := cfg.Get("cors").Populate(corsOpts); err != nil {
        panic(err)
    }

	msgOpts := controller.MessageOptions{
		Destination: sesOpts.To,
	}
	srv := createServer(*port, corsOpts.Origins, handlers{
		"POST /message": controller.NewMessage(client, validate, log, msgOpts),
	})

	done := launch(srv, log)
	select {
	case err = <-done:
		panic(err)
	case <-ctx.Done():
		break
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}

	if err := <-done; err != nil {
		panic(err)
	}
}

func quit() {
	if r := recover(); r != nil {
		fmt.Fprintf(os.Stderr, "unexpected panic occurred: %v\n", r)
		os.Exit(1)
	}
}
