package serverapp

import (
	"fmt"
	"log"
	"net/http"
)

type ServerApp struct {
	Server *http.Server
}

func New(serverAddress string, serverPort string, handler http.Handler) *ServerApp {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", serverAddress, serverPort),
		Handler: handler,
	}
	return &ServerApp{Server: server}
}

func (srv *ServerApp) MustRun() {
	if err := srv.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
