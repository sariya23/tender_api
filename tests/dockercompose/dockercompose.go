package dockercompose

import (
	"context"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/sariya23/tender/internal/config"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartComposeApp(ctx context.Context, pathToDockerCompose string, cfg *config.AppConfig) tc.ComposeStack {
	compose, err := tc.NewDockerCompose(pathToDockerCompose)
	if err != nil {
		panic(err)
	}
	composeWithEnvs := compose.WithEnv(map[string]string{
		"POSTGRES_DB":        cfg.PostgresDatabase,
		"POSTGRES_USERNAME":  cfg.PostgresUsername,
		"POSTRGRES_PASSWORD": cfg.PostgresPassword,
		"POSTGRES_PORT":      cfg.PostgresPassword,
		"SERVER_PORT":        cfg.ServerPort,
	})
	err = composeWithEnvs.Up(ctx, tc.Wait(false))
	if err != nil {
		panic(err)
	}
	composeWithEnvs.WaitForService("app",
		wait.ForHTTP("/api/ping").
			WithPort(nat.Port(cfg.ServerPort)).
			WithStartupTimeout(90*time.Second))
	return composeWithEnvs
}

func StopComposeApp(ctx context.Context, composeApp tc.ComposeStack) {
	err := composeApp.Down(ctx, tc.RemoveImagesAll, tc.RemoveOrphans(true), tc.RemoveImagesLocal)
	if err != nil {
		panic(err)
	}
}
