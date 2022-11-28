package handlers

import (
	"auth_service/application"
	"auth_service/domain"
	"auth_service/store"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type AuthHandler struct {
	service *application.AuthService
	store   *store.AuthMongoDBStore
}

func NewAuthHandler(service *application.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (handler *AuthHandler) Init(router *mux.Router) {
	router.HandleFunc("/login", handler.Login).Methods("POST")
	router.HandleFunc("/login", handler.Register).Methods("POST")
	http.Handle("/", router)
}

func (handler *AuthHandler) Register(writer http.ResponseWriter, req *http.Request) {
	var request domain.User
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	statusCode, err := handler.service.Register(&request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(statusCode)
}

func (handler *AuthHandler) Login(writer http.ResponseWriter, req *http.Request) {

	var request domain.User
	err := json.NewDecoder(req.Body).Decode(&request)

	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

}
