package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

const defaultTimeout = 10 * time.Second

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", defaultTimeout, "connection timeout")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	tClient := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := tClient.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, "Connection error:", err)
		os.Exit(1)
	}

	defer func() {
		err := tClient.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Connection close error:", err)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		if err := tClient.Send(); err != nil {
			fmt.Fprintln(os.Stderr, "Send error:", err)
			cancel()
		}
	}()

	go func() {
		if err := tClient.Receive(); err != nil {
			fmt.Fprintln(os.Stderr, "Receive error:", err)
			cancel()
		}
	}()

	<-ctx.Done()
}
