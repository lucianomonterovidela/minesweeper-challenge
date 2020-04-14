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

type IHandlerUser interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

type HandlerUser struct {
	userService services.IUserService
}

func (handler *HandlerUser) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user *models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	if err := validateUserRequest(user); err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
	}

	insertId, err := handler.userService.InsertUser(user)
	if err != nil {
		utils.LogError(err)
		server.InternalServerError(w, r, err)
		return
	}

	bodyResponse := make(map[string]interface{})
	bodyResponse["id"] = insertId
	server.OK(w, r, bodyResponse)
}

func (handler *HandlerUser) Login(w http.ResponseWriter, r *http.Request) {
	var user *models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
		return
	}

	if err := validateUserRequest(user); err != nil {
		server.BadRequest(w, r, server.ErrorCodeInvalidParams, err.Error())
	}

	token, err := handler.userService.Login(user.UserName, user.Password)

	if err != nil {
		server.InternalServerError(w, r, err)
	}

	if token == "" {
		server.Forbidden(w, r, "The user or password is incorrect.")
		return
	}

	bodyResponse := make(map[string]interface{})
	bodyResponse["token"] = token
	server.OK(w, r, bodyResponse)
}

func validateUserRequest(user *models.User) error {
	if user.UserName == "" {
		return fmt.Errorf("userName is mandatory")
	}

	if user.Password == "" {
		return fmt.Errorf("password is mandatory")
	}

	return nil
}

func NewHandlerUser() IHandlerUser {
	userService := services.NewUserService()
	return &HandlerUser{
		userService: userService,
	}
}
