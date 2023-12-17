package main

import (
	"context"
	"flag"
	"github.com/denismitr/antiddos/internal/bootstrap"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	host := flag.String("host", "127.0.0.1", "server host")
	port := flag.Int("port", 3333, "server port")
	zeroes := flag.Uint("zeroes", 3, "number of zeroes in hash")
	maxDuration := flag.Uint("max-duration", 30, "maximum duration of challenge in seconds")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		terminate := make(chan os.Signal, 1)
		signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)
		<-terminate
		cancel()
	}()

	c := bootstrap.TcpClient(uint8(*zeroes), uint64(*maxDuration), *host, *port)

	slog.Info("starting client")
	if err := c.Run(ctx); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info("client stopped")
}
