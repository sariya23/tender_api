package suite

import (
	"net/http"
	"testing"
	"time"
)

type Suite struct {
	T      *testing.T
	Client http.Client
}

func New(t *testing.T) *Suite {
	t.Helper()
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	return &Suite{
		T:      t,
		Client: client,
	}
}
