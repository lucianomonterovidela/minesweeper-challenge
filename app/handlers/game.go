package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/pedidosya/minesweeper-API/app/models"
	"github.com/pedidosya/minesweeper-API/app/server"
	"github.com/pedidosya/minesweeper-API/app/services"
	"github.com/pedidosya/minesweeper-API/utils"
	"net/http"
)

type IHandlerGame interface {
	NewGame(w http.ResponseWriter, r *http.Request)
	PauseGame(w http.ResponseWriter, r *http.Request)
	ResumeGame(w http.ResponseWriter, r *http.Request)
	Uncover(w http.ResponseWriter, r *http.Request)
	MarkRed(w http.ResponseWriter, r *http.Request)
	MarkQuestion(w http.ResponseWriter, r *http.Request)
	FindGames(w http.ResponseWriter, r *http.Request)
}

type HandlerGame struct {
	gameService services.IGameService
	userService services.IUserService
}

const authorizationHeader string = "Authorization"

func (handler *HandlerGame) NewGame(w http.ResponseWriter, r *http.Request) {
	var newGameRequest *models.NewGameRequest

	if err := json.NewDecoder(r.Body).Decode(&newGameRequest); err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	if err := validateNewGameRequest(newGameRequest); err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	userLogin := handler.userService.UserLogin(r.Header.Get(authorizationHeader))

	if userLogin == "" {
		server.Forbidden(w, r, "invalid token")
		return
	}

	game, err := handler.gameService.NewGame(newGameRequest.Rows, newGameRequest.Columns, newGameRequest.Mines, userLogin)
	if err != nil {
		server.InternalServerError(w, r, err)
		return
	}

	bodyResponse := make(map[string]interface{})
	bodyResponse["game"] = game
	server.OK(w, r, bodyResponse)
}

func (handler *HandlerGame) PauseGame(w http.ResponseWriter, r *http.Request) {

	gameId := server.GetStringFromPath(r, "game_id", "")
	if gameId == "" {
		err := fmt.Errorf("game id is mandatory")
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
	}

	userLogin := handler.userService.UserLogin(r.Header.Get(authorizationHeader))

	if userLogin == "" {
		server.Forbidden(w, r, "invalid token")
		return
	}

	isPaused, err := handler.gameService.PauseGame(gameId)
	if err != nil {
		utils.LogError(err)
		server.InternalServerError(w, r, err)
		return
	}

	if !isPaused {
		server.NotFound(w, r, fmt.Sprintf("Not found game wih gameId: %s", gameId))
		return
	}

	server.OkNotContent(w, r)
}

func (handler *HandlerGame) ResumeGame(w http.ResponseWriter, r *http.Request) {

	gameId := server.GetStringFromPath(r, "game_id", "")
	if gameId == "" {
		err := fmt.Errorf("game id is mandatory")
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
	}

	userLogin := handler.userService.UserLogin(r.Header.Get(authorizationHeader))

	if userLogin == "" {
		server.Forbidden(w, r, "invalid token")
		return
	}

	game, err := handler.gameService.ResumeGame(gameId, userLogin)
	if err != nil {
		utils.LogError(err)
		server.InternalServerError(w, r, err)
		return
	}

	bodyResponse := make(map[string]interface{})
	bodyResponse["game"] = game
	server.OK(w, r, bodyResponse)
}

func (handler *HandlerGame) MarkRed(w http.ResponseWriter, r *http.Request) {

	gameId := server.GetStringFromPath(r, "game_id", "")
	if gameId == "" {
		err := fmt.Errorf("game id is mandatory")
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
	}

	userLogin := handler.userService.UserLogin(r.Header.Get(authorizationHeader))

	if userLogin == "" {
		server.Forbidden(w, r, "invalid token")
		return
	}


	var cellRequest *models.CellRequest

	if err := json.NewDecoder(r.Body).Decode(&cellRequest); err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	if err := validateCellRequest(cellRequest); err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	isMark, err := handler.gameService.MarkRed(gameId, cellRequest.Row, cellRequest.Column)
	if err != nil {
		utils.LogError(err)
		server.InternalServerError(w, r, err)
		return
	}

	if !isMark {
		server.NotFound(w, r, fmt.Sprintf("not found game wih gameId: %s", gameId))
		return
	}

	server.OkNotContent(w, r)
}

func (handler *HandlerGame) MarkQuestion(w http.ResponseWriter, r *http.Request) {

	gameId := server.GetStringFromPath(r, "game_id", "")
	if gameId == "" {
		err := fmt.Errorf("game id is mandatory")
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
	}

	userLogin := handler.userService.UserLogin(r.Header.Get(authorizationHeader))

	if userLogin == "" {
		server.Forbidden(w, r, "invalid token")
		return
	}

	var cellRequest *models.CellRequest

	if err := json.NewDecoder(r.Body).Decode(&cellRequest); err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	if err := validateCellRequest(cellRequest); err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	isMark, err := handler.gameService.MarkQuestion(gameId, cellRequest.Row, cellRequest.Column)
	if err != nil {
		utils.LogError(err)
		server.InternalServerError(w, r, err)
		return
	}

	if !isMark {
		server.NotFound(w, r, fmt.Sprintf("not found game wih gameId: %s", gameId))
		return
	}

	server.OkNotContent(w, r)
}

func (handler *HandlerGame) Uncover(w http.ResponseWriter, r *http.Request) {

	gameId := server.GetStringFromPath(r, "game_id", "")
	if gameId == "" {
		err := fmt.Errorf("game id is mandatory")
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
	}

	userLogin := handler.userService.UserLogin(r.Header.Get(authorizationHeader))

	if userLogin == "" {
		server.Forbidden(w, r, "invalid token")
		return
	}

	var cellRequest *models.CellRequest

	if err := json.NewDecoder(r.Body).Decode(&cellRequest); err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	if err := validateCellRequest(cellRequest); err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	game, err := handler.gameService.Uncover(gameId, cellRequest.Row, cellRequest.Column)
	if err != nil {
		utils.LogError(err)
		server.InternalServerError(w, r, err)
		return
	}

	bodyResponse := make(map[string]interface{})
	bodyResponse["game"] = game
	server.OK(w, r, bodyResponse)
}

func (handler *HandlerGame) FindGames(w http.ResponseWriter, r *http.Request) {
	userLogin := handler.userService.UserLogin(r.Header.Get(authorizationHeader))

	if userLogin == "" {
		server.Forbidden(w, r, "invalid token")
		return
	}

	games, err := handler.gameService.FindGames(userLogin)

	if err == nil {
		server.OK(w, r, games)
	} else {
		server.InternalServerError(w, r, err)
	}
}

func validateCellRequest(cellRequest *models.CellRequest) error {
	if cellRequest.Row == 0 {
		return fmt.Errorf("row is mandatory and greater than zero")
	}

	if cellRequest.Column == 0 {
		return fmt.Errorf("columns is mandatory and greater than zero")
	}

	return nil
}

func validateNewGameRequest(newGameRequest *models.NewGameRequest) error {
	if newGameRequest.Rows == 0 {
		return fmt.Errorf("rows is mandatory and greater than zero")
	}

	if newGameRequest.Columns == 0 {
		return fmt.Errorf("columns is mandatory and greater than zero")
	}

	if newGameRequest.Mines == 0 {
		return fmt.Errorf("mines is mandatory and greater than zero")
	}

	if newGameRequest.Rows*newGameRequest.Columns <= newGameRequest.Mines {
		return fmt.Errorf("the board must have more cells than mines")
	}

	return nil
}

func NewHandlerGame() IHandlerGame {
	gameService := services.NewGameService()
	userService := services.NewUserService()
	return &HandlerGame{
		gameService: gameService,
		userService: userService,
	}
}
