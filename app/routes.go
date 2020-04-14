package main

import (
	"github.com/pedidosya/minesweeper-API/app/handlers"
	"github.com/pedidosya/minesweeper-API/app/server"
	"net/http"
)

func Routes(s *server.Server) {
	handlerRoot := handlers.NewHandlerHealth()
	s.AddRoute("/", handlerRoot.Health, http.MethodGet)

	handlerUser := handlers.NewHandlerUser()
	s.AddRoute("/v{version}/user", handlerUser.RegisterUser, http.MethodPost)
	s.AddRoute("/v{version}/user/login", handlerUser.Login, http.MethodPut)
}
