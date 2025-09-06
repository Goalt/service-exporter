package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Goalt/service-exporter/internal/app"
	"github.com/oklog/run"
)

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	app := app.New()
	if err := app.LoadConfig(); err != nil {
		log.Print("❌ Error loading config: ", err)
		return
	}

	var g run.Group
	{
		c := make(chan os.Signal, 1)

		g.Add(func() error {
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			s := <-c
			return fmt.Errorf("interrupted with sig %q", s)
		}, func(err error) {
			close(c)
			cancelCtx()
		})
	}
	{
		g.Add(func() error {
			if err := app.Run(ctx); err != nil {
				log.Print("❌ ", err)
			}

			return nil
		}, func(err error) {
			cancelCtx()
		})
	}

	if err := g.Run(); err != nil {
		log.Print("❌ Stopped with error: ", err)
	}
}
