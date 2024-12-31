package tests

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/sariya23/tender/internal/config"
	"github.com/sariya23/tender/internal/domain/models"
	schema "github.com/sariya23/tender/internal/hanlders"
	"github.com/sariya23/tender/internal/repository/postgres"
	"github.com/sariya23/tender/testdata"
	"github.com/sariya23/tender/tests/dockercompose"
	"github.com/sariya23/tender/tests/suite"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cfg := config.MustLoadByPath("../docker.env")
	defer cancel()
	app := dockercompose.StartComposeApp(ctx, "../docker-compose.yaml", cfg)
	db := postgres.MustNewConnection(ctx, cfg.PostgresConnOutside)
	err := db.CreateEmployee(ctx, testdata.TestEmployee)
	if err != nil {
		panic(err)
	}
	err = db.CreateOrganization(ctx, testdata.TestOrganization)
	if err != nil {
		panic(err)
	}
	tender, err := db.CreateTender(ctx, testdata.TestTender)
	if err != nil {
		panic(err)
	}
	_, err = db.EditTender(ctx, tender, 1, models.TenderToUpdate{Status: &models.TenderPublishedStatus})
	if err != nil {
		panic(err)
	}
	exitCode := m.Run()
	dockercompose.StopComposeApp(ctx, app)
	os.Exit(exitCode)
}

func TestGetAllTenders(t *testing.T) {
	st := suite.New(t)
	resp, err := st.Client.Get("http://127.0.0.1:8000/api/tenders")
	require.NoError(t, err)
	b := resp.Body
	defer b.Close()

	bytes, err := io.ReadAll(b)
	require.NoError(t, err)
	var respData schema.GetTendersResponse
	err = json.Unmarshal(bytes, &respData)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, schema.GetTendersResponse{Message: "ok", Tenders: []models.Tender{}}, respData)
}
