package handlers

import (
	"auth_service/application"
	"auth_service/domain"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type AuthHandler struct {
	service *application.AuthService
}

var jwtKey = []byte("my_secret_key")

func NewAuthHandler(service *application.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (handler *AuthHandler) Init(router *mux.Router) {
	router.HandleFunc("/login", handler.Login).Methods("POST")
	router.HandleFunc("/register", handler.Register).Methods("POST")
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

	err = handler.service.Register(&request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
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

	user, err1 := handler.service.Login(&request)

	if err1 != nil {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &domain.Claims{
		Username: user.Username,
		Role:     user.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(tokenString)

	http.SetCookie(writer, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})

}
