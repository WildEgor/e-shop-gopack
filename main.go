package main

import (
	slogger "github.com/WildEgor/e-shop-gopack/pkg/libs/logger/handlers"
	"log/slog"
	"os"
)

func main() {

	logger := slogger.NewLogger()
	slog.SetDefault(logger)

	slog.Info("Connect")
	os.Exit(1)
}
