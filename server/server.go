package server

import (
	"context"
	"net/http"
)

type HttpServer struct {
	server      *http.Server
	router      *http.ServeMux
	authHandler AuthUserHandler
}

func NewHttpServer(port string, authHandler AuthUserHandler) HttpServer {
	srv := HttpServer{
		router:      http.NewServeMux(),
		authHandler: authHandler,
	}
	srv.addRoutes()
	srv.server = &http.Server{Addr: port, Handler: srv.router}
	return srv
}

func (s HttpServer) Start() error {
	return s.server.ListenAndServe()
}

func (s HttpServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
