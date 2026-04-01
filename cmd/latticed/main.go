package main

import (
	"fmt"
	"log/slog"
	"os"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Printf("latticed %s\n", version)
	return nil
}
