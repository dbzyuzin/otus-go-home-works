package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

const minArgsCount = 3
const defaultTimeout = 10

func init() {
	flag.DurationVar(&timeout, "timeout", defaultTimeout*time.Second, "connection timeout")
}

func main() {
	argsCount := len(os.Args)
	if argsCount < minArgsCount {
		log.Fatalf("Expected to have at least 3 arguments, but got %d", argsCount)
	}
	flag.Parse()

	host, port := os.Args[2], os.Args[3]
	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = client.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	handlerWrapper := func(handler func() error, cancelFunc context.CancelFunc) {
		if err := handler(); err != nil {
			cancelFunc()
		}
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	go handlerWrapper(client.Receive, cancelFunc)
	go handlerWrapper(client.Send, cancelFunc)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigCh:
		cancelFunc()
		signal.Stop(sigCh)
		return

	case <-ctx.Done():
		signal.Stop(sigCh)
		close(sigCh)
		return
	}
}
