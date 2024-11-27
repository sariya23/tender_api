package serverapp

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"tender/internal/routes"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerApp struct {
	logger *slog.Logger
	server *http.Server
}

func New(logger *slog.Logger, addr string) *ServerApp {
	r := gin.Default()
	api := r.Group("api")
	{
		routes.BindRoutes(api)
		routes.TenderRoutes(api)
	}
	srv := &http.Server{Addr: addr, Handler: r}
	return &ServerApp{
		logger: logger,
		server: srv,
	}
}

func (s *ServerApp) MustRun() {
	if err := s.run(); err != nil {
		log.Fatal(err)
	}
}

func (s *ServerApp) run() error {
	const op = "app.server.run"
	logger := s.logger.With("op", op)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *ServerApp) GracefullStop(ctx context.Context) {
	const op = "app.server.GracefullStop"
	logger := s.logger.With("op", op)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown: ", slog.String("err", err.Error()))
		os.Exit(1)
	}
	select {
	case <-ctx.Done():
		logger.Info("timeout of 5 seconds.")
	}
	logger.Info("Server exiting")
}
