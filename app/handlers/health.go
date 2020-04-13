package handlers

import (
	"github.com/pedidosya/minesweeper-API/app/server"
	"github.com/spf13/viper"
	"net/http"
)

type IHandleHealth interface {
	Health(w http.ResponseWriter, r *http.Request)
}

type HandlerHealth struct{}

func (*HandlerHealth) Health(w http.ResponseWriter, r *http.Request) {
	server.OK(w, r, map[string]interface{}{
		"name": "minesweeper",
		"info": getInfo(),
	})
}

func getInfo() map[string]interface{} {
	return map[string]interface{}{
		"env":     viper.Get("server.port"),
		"version": viper.Get("server.version"),
	}

}

func NewHandlerHealth() IHandleHealth {
	return &HandlerHealth{}
}
