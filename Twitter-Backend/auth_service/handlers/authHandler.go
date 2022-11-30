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
	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/validateAccount", handler.ValidateAccount).Methods("POST")
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

	token, statusCode, err := handler.service.Register(&request)
	if err != nil {
		http.Error(writer, err.Error(), statusCode)
		return
	}

	jsonResponse(token, writer)
}

func (handler *AuthHandler) ValidateAccount(writer http.ResponseWriter, req *http.Request) {

	var request domain.RegisterValidation
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.service.ValidateAccount(&request)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusNotAcceptable)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (handler *AuthHandler) Login(writer http.ResponseWriter, req *http.Request) {

	var request domain.Credentials
	err := json.NewDecoder(req.Body).Decode(&request)

	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := handler.service.Login(&request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(token, writer)

}
