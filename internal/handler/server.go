package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"credit-plus/internal/service"
	"credit-plus/middleware"
)

type Server struct {
	port int

	service service.Service
}

func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8000
	}

	db := service.New()

	NewServer := &Server{
		port:    port,
		service: db,
	}

	basicAuthUsername := os.Getenv("BASIC_AUTH_USERNAME")
	basicAuthPassword := os.Getenv("BASIC_AUTH_PASSWORD")

	basicAuthMiddleware := middleware.NewBasicAuthMiddleware(middleware.BasicAuthConfig{
		Username: basicAuthUsername,
		Password: basicAuthPassword,
	})

	handler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      middleware.LoggerMiddleware(NewServer.RegisterRoutes(basicAuthMiddleware)),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
