package dockercompose

import (
	"context"
	"time"

	"github.com/sariya23/tender/testdata"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

func StartComposeApp(ctx context.Context, pathToDockerCompose string) tc.ComposeStack {
	compose, err := tc.NewDockerCompose(pathToDockerCompose)
	if err != nil {
		panic(err)
	}
	composeWithEnvs := compose.WithEnv(map[string]string{
		"POSTGRES_DB":        testdata.PostgresDBName,
		"POSTGRES_USERNAME":  testdata.PostgresUsername,
		"POSTRGRES_PASSWORD": testdata.PostgresPassword,
		"POSTGRES_PORT":      testdata.PostgresPort,
		"SERVER_PORT":        testdata.ServerPort,
	})
	err = composeWithEnvs.Up(ctx, tc.Wait(false))
	if err != nil {
		panic(err)
	}
	composeWithEnvs.WaitForService("app",
		wait.ForHTTP("/api/ping").
			WithPort(testdata.ServerPort).
			WithStartupTimeout(90*time.Second))
	return composeWithEnvs
}

func StopComposeApp(ctx context.Context, composeApp tc.ComposeStack) {
	err := composeApp.Down(ctx, tc.RemoveImagesAll)
	if err != nil {
		panic(err)
	}
}
