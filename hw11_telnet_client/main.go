package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func main() {
	if len(os.Args[1:]) < 2 {
		log.Fatalln("Not enough arguments")
	}

	flag.DurationVar(&timeout, "timeout", 0, "connection timeout")
	flag.Parse()

	params := os.Args[len(os.Args)-2:]
	address := params[0] + ":" + params[1]

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	if _, err := os.Stderr.WriteString("...Connected to " + address + "\n"); err != nil {
		cancel()
	} else {
		go process(ctx, cancel, client.Receive)
		go process(ctx, cancel, client.Send)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGQUIT, syscall.SIGINT)

	for {
		select {
		case sig := <-sigChan:
			if sig == syscall.SIGQUIT {
				_, _ = os.Stderr.WriteString("...EOF\n")
			}
			cancel()
		case <-ctx.Done():
			_ = client.Close()
			return
		}
	}
}

func process(ctx context.Context, cancel context.CancelFunc, processFunc func() error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := processFunc(); err != nil {
				_, _ = os.Stderr.WriteString(err.Error() + "\n")
				cancel()
			}
		}
	}
}
