package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Config struct {
	Port    int    `json:"port"`
	Version string `json:"version"`
}

type Server struct {
	httpSrv *http.Server
	router  *mux.Router
	cfg     *Config
	Version string
}

func New(c *Config) *Server {
	r := mux.NewRouter()
	return &Server{
		cfg: c,
		httpSrv: &http.Server{
			Handler:      r,
			Addr:         fmt.Sprintf(":%d", c.Port),
			WriteTimeout: 60 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
		router:  r,
		Version: c.Version,
	}
}

func (s *Server) AddRoute(path string, h http.HandlerFunc, methods ...string) {

	r := handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(h)

	g := http.TimeoutHandler(r, time.Duration(10)*time.Second, "response timeout exceeded")

	s.router.Handle(path, g).Methods(methods...)
}

func (s *Server) ListenAndServe() {
	s.httpSrv.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func Render(w http.ResponseWriter, r *http.Request, obj interface{}, status int) {
	js, err := json.Marshal(obj)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

func OK(w http.ResponseWriter, r *http.Request, obj interface{}) {
	Render(w, r, obj, http.StatusOK)
}

type errorResponse struct {
	Messages []string `json:"messages"`
	Code     string   `json:"code"`
}

func BadRequest(w http.ResponseWriter, r *http.Request, code string, messages ...string) {
	err := &errorResponse{
		Code:     code,
		Messages: messages,
	}
	Render(w, r, err, http.StatusBadRequest)
}

func Forbidden(w http.ResponseWriter, r *http.Request, messages ...string) {
	err := &errorResponse{
		Code:     "FORBIDDEN",
		Messages: messages,
	}
	Render(w, r, err, http.StatusForbidden)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	Render(w, r, &errorResponse{
		Code:     "INTERNAL_SERVER_ERROR",
		Messages: []string{err.Error()},
	}, http.StatusInternalServerError)
}