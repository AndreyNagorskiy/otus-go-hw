package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/AndreyNagorskiy/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	logger logger.Logger
	app    Application
	server *http.Server
}

type Application interface { // TODO
}

func NewServer(logger logger.Logger, app Application, host string, port int) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", hello)
	m := loggingMiddleware(mux)

	return &Server{
		logger: logger,
		app:    app,
		server: &http.Server{
			Addr:         host + ":" + strconv.Itoa(port),
			Handler:      m,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting server", slog.String("addr", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping server gracefully...")
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown failed: %w", err)
	}
	s.logger.Info("Server stopped")
	return nil
}

func hello(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("Hello, world!"))
}
