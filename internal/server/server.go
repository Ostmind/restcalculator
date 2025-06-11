package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"restcalculator/internal/http-server/middleware"
)

type Handler interface {
	Calculation(ctx echo.Context) error
}

type Server struct {
	handler Handler
	server  *echo.Echo
	logger  *slog.Logger
	port    int
}

func New(handler Handler, logger *slog.Logger, port int) *Server {

	server := echo.New()

	server.Use(middleware.LogRequest)
	server.POST("/:operation", handler.Calculation)

	return &Server{
		handler: handler,
		logger:  logger,
		server:  server,
		port:    port,
	}
}
func (s Server) Run() {
	s.logger.Info("Server is running on: localhost:%d", s.port)
	if err := s.server.Start(fmt.Sprintf("localhost:%d", s.port)); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Server starting error: %v", err)
		}
	}
}

func (s Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping server...")
	err := s.server.Shutdown(ctx)

	if err != nil {
		s.logger.Error("Error: ", err)
	}

	return err
}
