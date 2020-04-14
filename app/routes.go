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
	s.AddRoute("/v{version}/users", handlerUser.RegisterUser, http.MethodPost)
	s.AddRoute("/v{version}/users/login", handlerUser.Login, http.MethodPut)

	handlerGame := handlers.NewHandlerGame()
	s.AddRoute("/v{version}/games", handlerGame.NewGame, http.MethodPost)
	s.AddRoute("/v{version}/games/{game_id}/pause", handlerGame.PauseGame, http.MethodPut)
	s.AddRoute("/v{version}/games/{game_id}/resume", handlerGame.ResumeGame, http.MethodPut)
	s.AddRoute("/v{version}/games/{game_id}/mark-red", handlerGame.MarkRed, http.MethodPut)
	s.AddRoute("/v{version}/games/{game_id}/mark-question", handlerGame.MarkQuestion, http.MethodPut)
	s.AddRoute("/v{version}/games/{game_id}/uncover", handlerGame.Uncover, http.MethodPut)
	s.AddRoute("/v{version}/games", handlerGame.FindGames, http.MethodGet)
}
