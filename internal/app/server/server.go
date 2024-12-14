package serverapp

import (
	"log"
	"net/http"
)

type ServerApp struct {
	server *http.Server
}

func New(addr string, handler http.Handler) *ServerApp {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return &ServerApp{server: server}
}

func (srv *ServerApp) MustRun() {
	if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
