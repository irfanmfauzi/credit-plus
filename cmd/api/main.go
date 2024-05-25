package main

import (
	"credit-plus/internal/handler"
	"fmt"
	"log/slog"
)

func main() {

	server := handler.NewServer()

	slog.Info("Starting Credit Plus")

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
