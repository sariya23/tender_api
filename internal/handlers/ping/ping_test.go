package ping_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sariya23/tender/internal/handlers/ping"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPingHandler проверяет, что
// хендлер ping возвращает json {"message": "ok"} и статус ответа 200.
func TestPingHandler(t *testing.T) {
	r := gin.Default()
	r.GET("/ping", ping.Ping)
	req, _ := http.NewRequest("GET", "/ping", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	expectedBody := `{"message":"ok"}`
	require.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expectedBody, w.Body.String())
}
