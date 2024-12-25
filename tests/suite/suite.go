package suite

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/sariya23/tender/internal/config"
)

type Suite struct {
	T      *testing.T
	Client http.Client
	Config *config.AppConfig
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	cfg := config.MustLoadByPath("../local.env")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout))
	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	return ctx, &Suite{
		T:      t,
		Client: client,
		Config: cfg,
	}
}
