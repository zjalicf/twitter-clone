package handlers

import (
	"auth_service/application"
	"auth_service/domain"
	"auth_service/errors"
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
	router.HandleFunc("/verifyAccount", handler.VerifyAccount).Methods("POST")
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

func (handler *AuthHandler) VerifyAccount(writer http.ResponseWriter, req *http.Request) {

	var request domain.RegisterValidation
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if len(request.UserToken) == 0 {
		http.Error(writer, errors.InvalidUserTokenError, http.StatusBadRequest)
		return
	}

	err = handler.service.VerifyAccount(&request)
	if err != nil {
		if err.Error() == errors.InvalidTokenError {
			log.Println(err.Error())
			http.Error(writer, errors.InvalidTokenError, http.StatusNotAcceptable)
		} else if err.Error() == errors.ExpiredTokenError {
			log.Println(err.Error())
			http.Error(writer, errors.ExpiredTokenError, http.StatusNotFound)
		}
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
